package services

import (
	"github.com/fpay/foundation-go/log"
	"github.com/owenliu1122/notice"
)

// NewWeChatSenderService return a wechat sender service.
func NewWeChatSenderService(logger *log.Logger, toolCfg notice.SendServiceConfig, pc notice.ProducerInterface, exRouting notice.ProducerConfig) *WeChatSenderService {
	return &WeChatSenderService{
		logger: logger,
	}
}

// WeChatSenderService is a wechat sender service.
type WeChatSenderService struct {
	logger *log.Logger
}

// Handler parse a wechat message that needs to be sent.
func (svc *WeChatSenderService) Handler(msg *notice.UserMessage) error {

	// TODO: not implementation
	svc.logger.Debugf("WeChatSenderService: userMsg: %#v\n", msg)

	return nil
}
