package services

import (
	"fmt"

	"github.com/owenliu1122/notice"
)

// NewSenderService return target channel service interface.
func NewSenderService(msgType string, sendToolCfg notice.SendServiceConfig, pc notice.ProducerInterface, exRouting notice.ProducerConfig) SenderService {

	switch msgType {
	case NoticeTypeWeChat:
		return NewWeChatSenderService(sendToolCfg, pc, exRouting)
	case NoticeTypeMail:
		return NewMailSenderService(sendToolCfg, pc, exRouting)
	case NoticeTypePhone:
		return NewPhoneSenderService(sendToolCfg, pc, exRouting)
	default:
		panic(fmt.Sprintf("Unknown MsgType: %s", msgType))
	}
}

// SenderService sender message to target channel service interface.
type SenderService interface {
	Handler(msg *notice.UserMessage) error
}
