package controllers

import (
	"github.com/labstack/echo"
	log "gopkg.in/cihub/seelog.v2"
	"message_notification_practice/services"
	"net/http"
)

func NewGroupController(svc *services.GroupService) *GroupController {
	return &GroupController{svc: svc}
}

type GroupController struct {
	svc *services.GroupService
}

func (ctl *GroupController) List(ctx echo.Context) error {

	groups, err := ctl.svc.Find(0)

	if err != nil {
		log.Error("get groups list failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, groups)
}
