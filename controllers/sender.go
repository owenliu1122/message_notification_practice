package controllers

import (
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	log "gopkg.in/cihub/seelog.v2"

	"message_notification_practice/model"
	"message_notification_practice/services"
)

type SenderController struct {
	sendSvc services.SenderService
}

func NewSenderController(sendSvc services.SenderService) *SenderController {
	return &SenderController{
		sendSvc: sendSvc,
	}
}

func (ctl *SenderController) Handler(ctx context.Context, msg *amqp.Delivery) {

	userMsg := model.UserMsg{}

	if err := json.Unmarshal(msg.Body, &userMsg); err != nil {
		log.Error("Unmarshal MsgNotificationRequest Body failed, err: ", err)
		return
	}

	log.Debugf("SenderController:%T, %#v\n", ctl.sendSvc, userMsg)

	if err := ctl.sendSvc.Handler(&userMsg); err != nil {
		log.Error("get an error, handle it, err: ", err)
		return
	}

	return
}
