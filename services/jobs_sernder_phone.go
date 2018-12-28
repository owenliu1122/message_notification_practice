package services

import (
	log "gopkg.in/cihub/seelog.v2"
	"message_notification_practice"
)

func NewPhoneSenderService(toolCfg map[string]string) *PhoneSenderService {
	return &PhoneSenderService{}
}

type PhoneSenderService struct{}

func (svc *PhoneSenderService) Handler(msg *root.UserMsg) error {

	// TODO: not implementation
	log.Debugf("PhoneSenderService: userMsg: %#v\n", msg)

	return nil
}
