package services

import (
	"encoding/json"
	log "gopkg.in/cihub/seelog.v2"
	"message_notification_practice/model"
	"message_notification_practice/mq"
	"time"
)

const (
	MsgTypeMail   = "mail"
	MsgTypePhone  = "phone"
	MsgTypeWeChat = "wechat"
)

type MqSendService struct {
	mq        *mq.BaseMq
	gurSvc    *GroupUserRelationService
	exRouting map[string]*mq.BaseProducer
}

func NewMqSendService(mq *mq.BaseMq, gurSvc *GroupUserRelationService) *MqSendService {
	svc := MqSendService{
		mq:     mq,
		gurSvc: gurSvc,
	}

	return &svc
}

func (svc *MqSendService) RegisterExchangeRouting(tp string, exRouting mq.BaseProducer) {
	if svc.exRouting == nil {
		svc.exRouting = make(map[string]*mq.BaseProducer)
	}
	svc.exRouting[tp] = &exRouting
}

func (svc *MqSendService) Send(record *model.NotificationRecord) error {

	var err error
	var users []model.User
	users, err = svc.gurSvc.FindMembers(record.GroupID)
	if err != nil {
		log.Error("get group_user_relations failed, err: ", err)
	}
	for k, v := range svc.exRouting {

		log.Debugf("exRouting[%s]: %#v\n", k, v)
	}
	for _, user := range users {

		userMsg := &model.UserMsg{
			ID:      user.ID,
			Name:    user.Name,
			Content: record.Notification,
			Email:   user.Email,
			Phone:   user.Phone,
			WeChat:  user.Wechat,
		}

		body, err := json.Marshal(&userMsg)
		if err != nil {
			log.Error("Email marshal UserMsg failed, err: ", err)
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
