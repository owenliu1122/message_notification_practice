package controllers

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
	log "gopkg.in/cihub/seelog.v2"
	"message_notification_practice"
	"message_notification_practice/services"
)

type NotificationController struct {
	mqSvc *services.MqSendService
}

func NewNotificationController(mqSvc *services.MqSendService) *NotificationController {
	return &NotificationController{
		mqSvc: mqSvc,
	}
}

func (ctl *NotificationController) Handler(ctx context.Context, msg *amqp.Delivery) {
	var err error
	record := &root.NotificationRecord{}

	err = json.Unmarshal(msg.Body, record)
	if err != nil {
		log.Error("Unmarshal MsgNotificationRequest Body failed, err: ", err)
	}

	// TODO: return for test only
	log.Debugf("NotificationHandler -> record: %#v\n", record)

	err = ctl.mqSvc.Send(record)
	if err != nil {
		log.Error("mqSvc.Send record failed, err: ", err)
	}
}
