package services

import (
	"fmt"

	"github.com/owenliu1122/notice"
	"github.com/owenliu1122/notice/mq"
)

// NewSenderService return target channel service interface.
func NewSenderService(msgType string, sendToolCfg notice.SendService, pc *mq.Producer, exRouting notice.Producer) SenderService {

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
