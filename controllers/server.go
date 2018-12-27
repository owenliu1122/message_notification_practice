package controllers

import (
	"golang.org/x/net/context"
	log "gopkg.in/cihub/seelog.v2"
	"message_notification_practice/model"
	"message_notification_practice/pb"
	"message_notification_practice/services"
)

type ServerController struct {
	//mqChan    chan interface{}
	notifySrv *services.NotificationService
}

func NewServerController(notifySrv *services.NotificationService) *ServerController {
	return &ServerController{
		//mqChan:    mqChan,
		notifySrv: notifySrv,
	}
}

func (c *ServerController) CheckIn(ctx context.Context, request *pb.MsgNotificationRequest) (*pb.MsgNotificationResponse, error) {

	log.Debug(request.Content)

	if e := c.notifySrv.Create(&model.NotificationRecord{
		GroupID:      request.Group,
		Notification: request.Content,
	}); e != nil {
		log.Error("Insert record failed: ", e)
	}

	//c.mqChan <- jsonBytes

	return &pb.MsgNotificationResponse{Code: 0, Status: "success"}, nil
}
