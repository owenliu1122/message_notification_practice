package controllers

import (
	"context"
	"encoding/json"

	foundation "github.com/fpay/foundation-go"
	"github.com/fpay/foundation-go/log"
	"github.com/owenliu1122/notice"
	"github.com/owenliu1122/notice/services"
)

// SenderController is a sender notification controller
type SenderController struct {
	logger     *log.Logger
	sendSvc    services.SenderService
	retryDelay int
}

// NewSenderController returns a controller for sending notifications.
func NewSenderController(logger *log.Logger, retryDelay int, sendSvc services.SenderService) *SenderController {
	return &SenderController{
		logger:     logger,
		sendSvc:    sendSvc,
		retryDelay: retryDelay,
	}
}

// Handler parses the Sender controller
func (ctl *SenderController) Handler(ctx context.Context, job foundation.Jobber) error {
	var err error
	userMsg := notice.UserMessage{}

	if err = json.Unmarshal(job.Body(), &userMsg); err != nil {
		ctl.logger.Error("Unmarshal MsgNotificationRequest Body failed, err: ", err)
		return err
	}

	if err = ctl.sendSvc.Handler(ctx, &userMsg); err != nil {
		ctl.logger.Error("get an error, handle it, err: ", err.Error())
		return job.Retry(ctx, ctl.retryDelay)
	}

	return err
}
