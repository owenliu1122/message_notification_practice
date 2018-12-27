package controllers

import (
	//"database/sql/driver"
	"github.com/labstack/echo"
	log "gopkg.in/cihub/seelog.v2"
	"message_notification_practice"
	"message_notification_practice/services"
	"net/http"
	"strconv"
)

func NewGroupUserRelationController(svc *services.GroupUserRelationService) *GroupUserRelationController {
	return &GroupUserRelationController{svc: svc}
}

type GroupUserRelationController struct {
	svc *services.GroupUserRelationService
}

func (ctl *GroupUserRelationController) ListMembers(ctx echo.Context) error {

	groupStr := ctx.QueryParam("group_id")
	groupId, e := strconv.Atoi(groupStr)
	if e != nil {
		log.Errorf("group id string param convert to int, err: %s", groupStr, e)
		return ctx.String(http.StatusBadRequest, e.Error())
	}

	log.Debugf("GroupUserRelationController: groups_id: %d\n", groupId)

	users, err := ctl.svc.FindMembers(uint64(groupId))

	if err != nil {
		log.Errorf("get group(%d) members list failed, err: %s", groupId, err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, users)
}

func (ctl *GroupUserRelationController) AvailableMembers(ctx echo.Context) error {

	groupStr := ctx.QueryParam("group_id")
	searchUserName := ctx.QueryParam("user_name")
	groupId, e := strconv.Atoi(groupStr)
	if e != nil {
		log.Errorf("group id string param convert to int, err: %s", groupStr, e)
		return ctx.String(http.StatusBadRequest, e.Error())
	}

	log.Debugf("GroupUserRelationController: groups_id: %d, user_name: %s\n", groupId, searchUserName)

	users, err := ctl.svc.FindAvailableMembers(uint64(groupId), searchUserName)

	if err != nil {
		log.Errorf("get group(%d) members list failed, user_name: %s err: %s", groupId, searchUserName, err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, users)
}

func (ctl *GroupUserRelationController) Update(ctx echo.Context) error {

	var gur root.GroupUserRelation
	if err := ctx.Bind(&gur); err != nil {
		log.Error("update group get body failed, err: ", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	log.Infof("GroupController Update -> group: %#v\n", gur)

	err := ctl.svc.Update(&gur, nil)

	if err != nil {
		log.Error("update group failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, gur)
}

func (ctl *GroupUserRelationController) DeleteMembers(ctx echo.Context) error {

	var gur []root.GroupUserRelation

	if err := ctx.Bind(&gur); err != nil {
		log.Error("delete group user relations get body failed, err: ", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	if err := ctl.svc.Delete(gur); err != nil {
		log.Error("delete group user relations failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, gur)
}

func (ctl *GroupUserRelationController) AddMembers(ctx echo.Context) error {

	log.Debug("start AddMembers")

	var gur []root.GroupUserRelation

	if err := ctx.Bind(&gur); err != nil {
		log.Error("create group user relations get body failed, err: ", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	log.Debugf("AddMembers: %#v\n", gur)

	if err := ctl.svc.Create(gur); err != nil {
		log.Error("create group user relations failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, gur)
}
