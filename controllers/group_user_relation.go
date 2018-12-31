package controllers

import (
	//"database/sql/driver"
	"message_notification_practice"
	"message_notification_practice/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "gopkg.in/cihub/seelog.v2"
)

// NewGroupUserRelationController returns a controller for group user relations table.
func NewGroupUserRelationController(svc *services.GroupUserRelationService) *GroupUserRelationController {
	return &GroupUserRelationController{svc: svc}
}

// GroupUserRelationController is a group user relation controller
type GroupUserRelationController struct {
	svc *services.GroupUserRelationService
}

// ListMembers will return all members for current groups id.
func (ctl *GroupUserRelationController) ListMembers(ctx echo.Context) error {

	groupStr := ctx.QueryParam("group_id")
	groupID, e := strconv.Atoi(groupStr)
	if e != nil {
		log.Errorf("group id string param convert to int, err: %s", groupStr, e)
		return ctx.String(http.StatusBadRequest, e.Error())
	}

	log.Debugf("GroupUserRelationController: groups_id: %d\n", groupID)

	users, err := ctl.svc.FindMembers(uint64(groupID))
	if err != nil {
		log.Errorf("get group(%d) members list failed, err: %s", groupID, err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, users)
}

// AvailableMembers will list all users that can be add to current group id.
func (ctl *GroupUserRelationController) AvailableMembers(ctx echo.Context) error {

	groupStr := ctx.QueryParam("group_id")
	searchUserName := ctx.QueryParam("user_name")
	groupID, e := strconv.Atoi(groupStr)
	if e != nil {
		log.Errorf("group id string param convert to int, err: %s", groupStr, e)
		return ctx.String(http.StatusBadRequest, e.Error())
	}

	log.Debugf("GroupUserRelationController: groups_id: %d, user_name: %s\n", groupID, searchUserName)

	users, err := ctl.svc.FindAvailableMembers(uint64(groupID), searchUserName)

	if err != nil {
		log.Errorf("get group(%d) members list failed, user_name: %s err: %s", groupID, searchUserName, err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, users)
}

// Update group user relation record.
func (ctl *GroupUserRelationController) Update(ctx echo.Context) error {

	var gur notice.GroupUserRelation
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

// DeleteMembers pare group user relations deleting operations.
func (ctl *GroupUserRelationController) DeleteMembers(ctx echo.Context) error {

	var gur []notice.GroupUserRelation

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

// AddMembers parse group user reations creating operation.
func (ctl *GroupUserRelationController) AddMembers(ctx echo.Context) error {

	log.Debug("start AddMembers")

	var gur []notice.GroupUserRelation

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
