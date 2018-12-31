package controllers

import (
	"message_notification_practice"
	"message_notification_practice/services"
	"net/http"

	"github.com/labstack/echo"
	log "gopkg.in/cihub/seelog.v2"
)

// NewGroupController will return a groups table operation controller.
func NewGroupController(svc *services.GroupService) *GroupController {
	return &GroupController{svc: svc}
}

// GroupController is a groups table operation controller.
type GroupController struct {
	svc *services.GroupService
}

// List all group user relation records.
func (ctl *GroupController) List(ctx echo.Context) error {

	groups, err := ctl.svc.Find(0)

	if err != nil {
		log.Error("get groups list failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, groups)
}

// Create parse the group user relations table creating operations.
func (ctl *GroupController) Create(ctx echo.Context) error {

	var group notice.Group
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

// Update parse the group user relations table updating operations.
func (ctl *GroupController) Update(ctx echo.Context) error {

	var group notice.Group
	if err := ctx.Bind(&group); err != nil {
		log.Error("update group get body failed, err: ", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	log.Infof("GroupController Update -> group: %#v\n", group)

	err := ctl.svc.Update(&group)

	if err != nil {
		log.Error("update group failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, group)
}

// Delete parse the group user relations table deleting operations.
func (ctl *GroupController) Delete(ctx echo.Context) error {

	var group notice.Group
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
