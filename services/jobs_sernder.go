package services

import (
	"fmt"
	"message_notification_practice"
)

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

	return nil
}

type SenderService interface {
	Handler(msg *root.UserMsg) error
}
