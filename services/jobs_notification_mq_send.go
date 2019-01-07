package services

import (
	"context"
	"fmt"

	"github.com/fpay/foundation-go"

	"github.com/owenliu1122/notice"
	"github.com/owenliu1122/notice/pb"

	"github.com/fpay/foundation-go/log"
)

// send notification channel type
const (
	NoticeTypeMail   = "mail"
	NoticeTypePhone  = "phone"
	NoticeTypeWeChat = "wechat"
)

// MqSendService is mq send service.
type MqSendService struct {
	logger   *log.Logger
	grpSvc   *GroupService
	pc       foundation.JobManager
	queueCfg map[string]notice.JobConfig
}

// NewMqSendService returns a mq send service.
func NewMqSendService(logger *log.Logger, pc foundation.JobManager, grpSvc *GroupService, jobCfg map[string]notice.JobConfig) *MqSendService {
	svc := MqSendService{
		logger:   logger,
		pc:       pc,
		grpSvc:   grpSvc,
		queueCfg: jobCfg,
	}

	return &svc
}

// Send parse send a record to  exchange and routingkey.
func (svc *MqSendService) Send(ctx context.Context, record *pb.MsgNotificationRequest) error {

	var err error
	var users []notice.User
	users, err = svc.grpSvc.FindMembers(record.Group)
	if err != nil {
		return fmt.Errorf("get group_user_relations failed, err: %s", err)
	}

	for _, user := range users {

		userMsg := &notice.UserMessage{
			ID:      user.ID,
			Name:    user.Name,
			Content: record.Content,
		}

		for _, noticeType := range record.NoticeType {

			strType := noticeType.String()

			userMsg.NoticeType = strType

			switch strType {
			case NoticeTypeMail:
				userMsg.Destination = user.Email
			case NoticeTypePhone:
				userMsg.Destination = user.Phone
			case NoticeTypeWeChat:
				userMsg.Destination = user.Wechat
			default:
				return fmt.Errorf("unknown notice type: %s", strType)
			}

			queueCfg, ok := svc.queueCfg[strType]
			if !ok {
				return fmt.Errorf("get producer config failed, Unknown notice type: %s", strType)
			}

			if err = svc.pc.Dispatch(ctx, &Job{
				Q:       queueCfg.Queue,
				D:       queueCfg.Delay,
				Message: &userMsg,
			}); err != nil {
				return err
			}
		}
	}

	svc.logger.Debugf("group_id: %d, %#v\n", record.Group, users)

	return err
}
