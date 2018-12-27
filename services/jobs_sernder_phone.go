package services

import (
	log "gopkg.in/cihub/seelog.v2"
	"gopkg.in/mailgun/mailgun-go.v1"
	"message_notification_practice/model"
)

func NewPhoneSenderService(mg interface{}) *PhoneSenderService {
	return &PhoneSenderService{mg: mg.(mailgun.Mailgun)}
}

type PhoneSenderService struct {
	mg mailgun.Mailgun
}

func (svc *PhoneSenderService) Handler(msg *model.UserMsg) error {

	// TODO: not implementation
	log.Debugf("PhoneSenderService: userMsg: %#v\n", msg)

	return nil
}
