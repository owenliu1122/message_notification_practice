package services

import (
	"encoding/json"
	"message_notification_practice"
	"message_notification_practice/mq"
	"time"

	log "gopkg.in/cihub/seelog.v2"
)

// send notification channel type
const (
	MsgTypeMail   = "mail"
	MsgTypePhone  = "phone"
	MsgTypeWeChat = "wechat"
)

// MqSendService is mq send service.
type MqSendService struct {
	mq        *mq.BaseMq
	gurSvc    *GroupUserRelationService
	exRouting map[string]*mq.BaseProducer
}

// NewMqSendService returns a mq send service.
func NewMqSendService(mq *mq.BaseMq, gurSvc *GroupUserRelationService) *MqSendService {
	svc := MqSendService{
		mq:     mq,
		gurSvc: gurSvc,
	}

	return &svc
}

// RegisterExchangeRouting regist exchange and routingkey.
func (svc *MqSendService) RegisterExchangeRouting(tp string, exRouting mq.BaseProducer) {
	if svc.exRouting == nil {
		svc.exRouting = make(map[string]*mq.BaseProducer)
	}
	svc.exRouting[tp] = &exRouting
}

// Send parse send a record to  exchange and routingkey.
func (svc *MqSendService) Send(record *notice.NotificationRecord) error {

	var err error
	var users []notice.User
	users, err = svc.gurSvc.FindMembers(record.GroupID)
	if err != nil {
		log.Error("get group_user_relations failed, err: ", err)
	}
	for k, v := range svc.exRouting {
		log.Debugf("exRouting[%s]: %#v\n", k, v)
	}
	for _, user := range users {

		userMsg := &notice.UserMessage{
			ID:      user.ID,
			Name:    user.Name,
			Content: record.Notification,
			Email:   user.Email,
			Phone:   user.Phone,
			WeChat:  user.Wechat,
		}

		body, e := json.Marshal(&userMsg)
		if e != nil {
			log.Error("Email marshal UserMsg failed, err: ", e)
		}

		if len(user.Email) > 0 {
			if err = svc.mq.Send(svc.exRouting[MsgTypeMail].Exchange,
				svc.exRouting[MsgTypeMail].RoutingKey,
				body); err != nil {
				return err
			}
		}

		if len(user.Phone) > 0 {
			if err = svc.mq.Send(svc.exRouting[MsgTypePhone].Exchange,
				svc.exRouting[MsgTypePhone].RoutingKey,
				body); err != nil {
				return err
			}
		}

		if len(user.Wechat) > 0 {
			if err = svc.mq.Send(svc.exRouting[MsgTypeWeChat].Exchange,
				svc.exRouting[MsgTypeWeChat].RoutingKey,
				body); err != nil {
				return err
			}
		}
	}

	log.Debugf("group_id: %d, %#v\n", record.GroupID, users)

	time.Sleep(2 * time.Second) // TODO: remove debug

	return err
}
