package services

import (
	"fmt"
	"message_notification_practice/model"
)

func NewSenderService(msgType string, mg interface{}) SenderService {

	switch msgType {
	case MsgTypeWeChat:
		return NewWeChatSenderService(mg)
	case MsgTypeMail:
		return NewMailSenderService(mg)
	case MsgTypePhone:
		return NewPhoneSenderService(mg)
	default:
		panic(fmt.Sprintf("Unknown MsgType: %s", msgType))
	}

	return nil
}

type SenderService interface {
	Handler(msg *model.UserMsg) error
}
