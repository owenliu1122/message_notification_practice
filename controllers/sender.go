package controllers

import (
	"context"
	"encoding/json"

	"github.com/streadway/amqp"
	log "gopkg.in/cihub/seelog.v2"

	"github.com/owenliu1122/notice"
	"github.com/owenliu1122/notice/services"
)

// SenderController is a sender notification controller
type SenderController struct {
	sendSvc services.SenderService
}

// NewSenderController returns a controller for sending notifications.
func NewSenderController(sendSvc services.SenderService) *SenderController {
	return &SenderController{
		sendSvc: sendSvc,
	}
}

// Handler parses the Sender controller
func (ctl *SenderController) Handler(ctx context.Context, msg *amqp.Delivery) {

	userMsg := notice.UserMessage{}

	if err := json.Unmarshal(msg.Body, &userMsg); err != nil {
		log.Error("Unmarshal MsgNotificationRequest Body failed, err: ", err)
		return
	}

	log.Debugf("SenderController:%T, %#v\n", ctl.sendSvc, userMsg)

	if err := ctl.sendSvc.Handler(&userMsg); err != nil {
		log.Error("get an error, handle it, err: ", err.Error())
		return
	}

	return
}
