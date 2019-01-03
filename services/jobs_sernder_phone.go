package services

import (
	"github.com/owenliu1122/notice"
	log "gopkg.in/cihub/seelog.v2"
)

// NewPhoneSenderService return a mail sender service.
func NewPhoneSenderService(toolCfg notice.SendServiceConfig, pc notice.ProducerInterface, exRouting notice.ProducerConfig) *PhoneSenderService {
	return &PhoneSenderService{}
}

// PhoneSenderService is a phone sender service.
type PhoneSenderService struct{}

// Handler parse a phone message that needs to be sent.
func (svc *PhoneSenderService) Handler(msg *notice.UserMessage) error {

	// TODO: not implementation
	log.Debugf("PhoneSenderService: userMsg: %#v\n", msg)

	return nil
}
