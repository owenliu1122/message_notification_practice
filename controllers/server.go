package controllers

import (
	"github.com/owenliu1122/notice/pb"
	"github.com/owenliu1122/notice/services"

	"github.com/fpay/foundation-go/log"
	"golang.org/x/net/context"
)

// ServerController is a server controller.
type ServerController struct {
	//mqChan    chan interface{}
	logger    *log.Logger
	notifySrv *services.NotificationService
}

// NewServerController will returns a server controller.
func NewServerController(logger *log.Logger, notifySrv *services.NotificationService) *ServerController {
	return &ServerController{
		//mqChan:    mqChan,
		notifySrv: notifySrv,
	}
}

// CheckIn parse grpc client notification request.
func (c *ServerController) CheckIn(ctx context.Context, request *pb.MsgNotificationRequest) (*pb.MsgNotificationResponse, error) {

	c.logger.Debug(request.Content)

	if e := c.notifySrv.Create(request); e != nil {
		c.logger.Error("Insert record failed: ", e.Error())
	}

	return &pb.MsgNotificationResponse{Code: 0, Status: "success"}, nil
}
