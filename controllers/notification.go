package controllers

import (
	"encoding/json"
	"time"

	"github.com/fpay/foundation-go"

	"github.com/owenliu1122/notice/pb"
	"github.com/owenliu1122/notice/services"

	"github.com/fpay/foundation-go/log"
	"golang.org/x/net/context"
)

// NotificationController is a notification controller
type NotificationController struct {
	logger *log.Logger
	mqSvc  *services.MqSendService
}

// NewNotificationController returns a controller for service message notification and save.
func NewNotificationController(logger *log.Logger, mqSvc *services.MqSendService) *NotificationController {
	return &NotificationController{
		logger: logger,
		mqSvc:  mqSvc,
	}
}

// Handler parse rabbitmq notifications.
func (ctl *NotificationController) Handler(ctx context.Context, job foundation.Jobber) (err error) {

	record := &pb.MsgNotificationRequest{}

	err = json.Unmarshal(job.Body(), record)
	if err != nil {
		ctl.logger.Error("Unmarshal MsgNotificationRequest Body failed, err: ", err)
		return
	}

	err = ctl.mqSvc.Send(ctx, record)
	if err != nil {
		ctl.logger.Error("mqSvc.Send record failed, err: ", err)
		return
	}

	time.Sleep(5 * time.Second)

	return
}
