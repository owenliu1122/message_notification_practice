package controllers

import (
	"message_notification_practice/model"
)

type SenderController struct{}

type SenderHandler func(msg *model.UserMsg) error

func NewSenderController(tp string, params ...interface{}) SenderControllerInterface {
	//return NewMailSenderController(params[0].(string), params[0].(string), params[0].(string))
	return NewMailSenderController(params[0].(string), params[1].(string), params[2].(string))
}

type SenderControllerInterface interface {
	Handler(msg *model.UserMsg) error
}
