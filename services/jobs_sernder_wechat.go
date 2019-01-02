package services

import (
	"github.com/owenliu1122/notice"
	"github.com/owenliu1122/notice/mq"
	log "gopkg.in/cihub/seelog.v2"
)

// NewWeChatSenderService return a wechat sender service.
func NewWeChatSenderService(toolCfg notice.SendService, pc *mq.Producer, exRouting notice.Producer) *WeChatSenderService {
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
