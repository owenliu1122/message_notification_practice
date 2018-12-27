package services

import (
	log "gopkg.in/cihub/seelog.v2"
	"message_notification_practice"
)

func NewWeChatSenderService(toolCfg interface{}) *WeChatSenderService {
	return &WeChatSenderService{}
}

type WeChatSenderService struct{}

func (svc *WeChatSenderService) Handler(msg *root.UserMsg) error {

	// TODO: not implementation
	log.Debugf("WeChatSenderService: userMsg: %#v\n", msg)

	return nil
}
