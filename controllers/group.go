package controllers

import (
	"github.com/labstack/echo"
	log "gopkg.in/cihub/seelog.v2"
	"message_notification_practice/model"
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

func (ctl *GroupController) Create(ctx echo.Context) error {

	var group model.Group
	if err := ctx.Bind(&group); err != nil {
		log.Error("add group get body failed, err: ", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	log.Infof("GroupController Bind -> group: %v\n", group)

	err := ctl.svc.Create(&group)

	if err != nil {
		log.Error("create group failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, group)
}

func (ctl *GroupController) Update(ctx echo.Context) error {

	var group model.Group
	if err := ctx.Bind(&group); err != nil {
		log.Error("update group get body failed, err: ", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	log.Infof("GroupController Update -> group: %#v\n", group)

	err := ctl.svc.Update(&group, nil)

	if err != nil {
		log.Error("update group failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, group)
}

func (ctl *GroupController) Delete(ctx echo.Context) error {

	var group model.Group
	if err := ctx.Bind(&group); err != nil {
		log.Error("delete group get body failed, err: ", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	log.Infof("GroupController Delete -> group: %#v\n", group)

	_, err := ctl.svc.Delete(&group)

	if err != nil {
		log.Error("delete group failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, group)
}
