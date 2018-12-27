package services

import (
	log "gopkg.in/cihub/seelog.v2"
	"gopkg.in/mailgun/mailgun-go.v1"
	"message_notification_practice/model"
)

func NewWeChatSenderService(mg interface{}) *WeChatSenderService {
	return &WeChatSenderService{mg: mg.(mailgun.Mailgun)}
}

type WeChatSenderService struct {
	mg mailgun.Mailgun
}

func (svc *WeChatSenderService) Handler(msg *model.UserMsg) error {

	// TODO: not implementation
	log.Debugf("WeChatSenderService: userMsg: %#v\n", msg)

	return nil
}
