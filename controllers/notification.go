package controllers

import (
	"encoding/json"

	"github.com/owenliu1122/notice/pb"
	"github.com/owenliu1122/notice/services"

	"github.com/streadway/amqp"
	"golang.org/x/net/context"
	log "gopkg.in/cihub/seelog.v2"
)

// NotificationController is a notification controller
type NotificationController struct {
	mqSvc *services.MqSendService
}

// NewNotificationController returns a controller for service message notification and save.
func NewNotificationController(mqSvc *services.MqSendService) *NotificationController {
	return &NotificationController{
		mqSvc: mqSvc,
	}
}

// Handler parse rabbitmq notifications.
func (ctl *NotificationController) Handler(ctx context.Context, msg *amqp.Delivery) {
	var err error
	record := &pb.MsgNotificationRequest{}

	err = json.Unmarshal(msg.Body, record)
	if err != nil {
		log.Error("Unmarshal MsgNotificationRequest Body failed, err: ", err)
		return
	}

	err = ctl.mqSvc.Send(record)
	if err != nil {
		log.Error("mqSvc.Send record failed, err: ", err)
	}
}
