package services

import (
	"context"

	"github.com/fpay/foundation-go"

	"github.com/fpay/foundation-go/log"
	"github.com/owenliu1122/notice"
)

// NewPhoneSenderService return a mail sender service.
func NewPhoneSenderService(logger *log.Logger, toolCfg notice.SendServiceConfig, pc foundation.JobManager) *PhoneSenderService {
	return &PhoneSenderService{
		logger: logger,
	}
}

// PhoneSenderService is a phone sender service.
type PhoneSenderService struct {
	logger *log.Logger
}

// Handler parse a phone message that needs to be sent.
func (svc *PhoneSenderService) Handler(ctx context.Context, msg *notice.UserMessage) error {

	// TODO: not implementation
	svc.logger.Debugf("PhoneSenderService: userMsg: %#v\n", msg)

	return nil
}
