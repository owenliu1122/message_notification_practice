package services

import (
	"message_notification_practice"

	log "gopkg.in/cihub/seelog.v2"
)

// NewWeChatSenderService return a wechat sender service.
func NewWeChatSenderService(toolCfg notice.SendService) *WeChatSenderService {
	return &WeChatSenderService{}
}

// WeChatSenderService is a wechat sender service.
type WeChatSenderService struct{}

// Handler parse a wechat message that needs to be sent.
func (svc *WeChatSenderService) Handler(msg *notice.UserMessage) error {

	// TODO: not implementation
	log.Debugf("WeChatSenderService: userMsg: %#v\n", msg)

	return nil
}
