package services

import (
	"fmt"
	"message_notification_practice"
)

// NewSenderService return target channel service interface.
func NewSenderService(msgType string, sendToolCfg map[string]string) SenderService {

	switch msgType {
	case MsgTypeWeChat:
		return NewWeChatSenderService(sendToolCfg)
	case MsgTypeMail:
		return NewMailSenderService(sendToolCfg)
	case MsgTypePhone:
		return NewPhoneSenderService(sendToolCfg)
	default:
		panic(fmt.Sprintf("Unknown MsgType: %s", msgType))
	}
}

// SenderService sender message to target channel service interface.
type SenderService interface {
	Handler(msg *notice.UserMessage) error
}
