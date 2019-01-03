package controllers

import (
	"github.com/owenliu1122/notice/pb"
	"github.com/owenliu1122/notice/services"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// ServerController is a server controller.
type ServerController struct {
	//mqChan    chan interface{}
	notifySrv *services.NotificationService
}

// NewServerController will returns a server controller.
func NewServerController(notifySrv *services.NotificationService) *ServerController {
	return &ServerController{
		//mqChan:    mqChan,
		notifySrv: notifySrv,
	}
}

// CheckIn parse grpc client notification request.
func (c *ServerController) CheckIn(ctx context.Context, request *pb.MsgNotificationRequest) (*pb.MsgNotificationResponse, error) {

	log.Debug(request.Content)

	if e := c.notifySrv.Create(request); e != nil {
		log.Error("Insert record failed: ", e.Error())
	}

	//c.mqChan <- jsonBytes

	return &pb.MsgNotificationResponse{Code: 0, Status: "success"}, nil
}
