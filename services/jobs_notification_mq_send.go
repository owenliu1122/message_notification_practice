package services

import (
	"encoding/json"
	"fmt"

	"github.com/owenliu1122/notice"
	"github.com/owenliu1122/notice/pb"

	log "gopkg.in/cihub/seelog.v2"
)

// send notification channel type
const (
	NoticeTypeMail   = "mail"
	NoticeTypePhone  = "phone"
	NoticeTypeWeChat = "wechat"
)

// MqSendService is mq send service.
type MqSendService struct {
	pc        notice.ProducerInterface
	grpSvc    *GroupService
	exRouting map[string]notice.ProducerConfig
}

// NewMqSendService returns a mq send service.
func NewMqSendService(pc notice.ProducerInterface, grpSvc *GroupService, exRouting map[string]notice.ProducerConfig) *MqSendService {
	svc := MqSendService{
		pc:     pc,
		grpSvc: grpSvc,
	}
	svc.exRouting = make(map[string]notice.ProducerConfig)
	svc.exRouting = exRouting
	return &svc
}

// Send parse send a record to  exchange and routingkey.
func (svc *MqSendService) Send(record *pb.MsgNotificationRequest) error {

	var err error
	var users []notice.User
	users, err = svc.grpSvc.FindMembers(record.Group)
	if err != nil {
		return fmt.Errorf("get group_user_relations failed, err: %s", err)
	}

	for k, v := range svc.exRouting {
		log.Debugf("exRouting[%s]: %#v\n", k, v)
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

			producerCfg, ok := svc.exRouting[strType]
			if !ok {
				return fmt.Errorf("get producer config failed, Unknown notice type: %s", strType)
			}

			body, e := json.Marshal(&userMsg)
			if e != nil {
				return fmt.Errorf("email marshal UserMsg failed, err: %s", e)
			}

			if err = svc.pc.Publish(producerCfg.Exchange, producerCfg.RoutingKey, body); err != nil {
				return err
			}
		}
	}

	log.Debugf("group_id: %d, %#v\n", record.Group, users)

	return err
}
