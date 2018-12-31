package services

import (
	"fmt"
	"message_notification_practice"
)

// NewSenderService return target channel service interface.
func NewSenderService(msgType string, sendToolCfg notice.SendService) SenderService {

	switch msgType {
	case NoticeTypeWeChat:
		return NewWeChatSenderService(sendToolCfg)
	case NoticeTypeMail:
		return NewMailSenderService(sendToolCfg)
	case NoticeTypePhone:
		return NewPhoneSenderService(sendToolCfg)
	default:
		panic(fmt.Sprintf("Unknown MsgType: %s", msgType))
	}
}

// SenderService sender message to target channel service interface.
type SenderService interface {
	Handler(msg *notice.UserMessage) error
}
