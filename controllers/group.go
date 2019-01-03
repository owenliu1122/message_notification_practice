package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/owenliu1122/notice"
	log "github.com/sirupsen/logrus"
)

// NewGroupController will return a groups table operation controller.
func NewGroupController(svc notice.GroupServiceInterface) *GroupController {
	return &GroupController{svc: svc}
}

// GroupController is a groups table operation controller.
type GroupController struct {
	svc notice.GroupServiceInterface
}

// List all group user relation records.
func (ctl *GroupController) List(ctx echo.Context) error {

	groupName := ctx.QueryParam("name")

	pageStr := ctx.QueryParam("page")
	page, e := strconv.Atoi(pageStr)
	if e != nil {
		log.Errorf("page string param convert to int, pagt: %s, err: %s", pageStr, e)
		return ctx.String(http.StatusBadRequest, e.Error())
	}

	pageSizeStr := ctx.QueryParam("page_size")
	pageSize, e := strconv.Atoi(pageSizeStr)
	if e != nil {
		log.Errorf("page size string param convert to int, page size: %s, err: %s", pageSizeStr, e)
		return ctx.String(http.StatusBadRequest, e.Error())
	}

	groups, cnt, err := ctl.svc.List(groupName, page, pageSize)

	if err != nil {
		log.Error("get groups list failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"count": cnt, "data": groups})
}

// Create parse the group user relations table creating operations.
func (ctl *GroupController) Create(ctx echo.Context) error {

	var group notice.Group
	if err := ctx.Bind(&group); err != nil {
		log.Error("add group get body failed, err: ", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	log.Infof("GroupController Bind -> group: %v\n", group)

	if group.Name == "" {
		err := errors.New("create group failed, err: no group name")
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err := ctl.svc.Create(&group)

	if err != nil {
		err = fmt.Errorf("create group failed, err: %s", err)
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

	if group.ID == 0 {
		err := errors.New("update group failed, err: no group id")
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	if group.Name == "" {
		err := errors.New("update group failed, err: no group name")
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err := ctl.svc.Update(&group)

	if err != nil {
		err = fmt.Errorf("update group failed, err: %s", err)
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

	if group.ID == 0 {
		err := errors.New("delete group failed, err: no group id")
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	_, err := ctl.svc.Delete(&group)

	if err != nil {
		log.Error("delete group failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, group)
}

// ListMembers will return all members for current groups id.
func (ctl *GroupController) ListMembers(ctx echo.Context) error {

	groupStr := ctx.QueryParam("group_id")
	groupID, e := strconv.Atoi(groupStr)
	if e != nil {
		log.Errorf("group id string param convert to int, group id: %s, err: %s", groupStr, e)
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
func (ctl *GroupController) AvailableMembers(ctx echo.Context) error {

	groupStr := ctx.QueryParam("group_id")
	searchUserName := ctx.QueryParam("user_name")
	groupID, e := strconv.Atoi(groupStr)
	if e != nil {
		log.Errorf("group id string param convert to int, group id: %s, err: %s", groupStr, e)
		return ctx.String(http.StatusBadRequest, e.Error())
	}

	pageStr := ctx.QueryParam("page")
	page, e := strconv.Atoi(pageStr)
	if e != nil {
		log.Errorf("page string param convert to int, page: %s, err: %s", pageStr, e)
		return ctx.String(http.StatusBadRequest, e.Error())
	}

	pageSizeStr := ctx.QueryParam("page_size")
	pageSize, e := strconv.Atoi(pageSizeStr)
	if e != nil {
		log.Errorf("page size string param convert to int, page size: %s, err: %s", pageSizeStr, e)
		return ctx.String(http.StatusBadRequest, e.Error())
	}

	log.Debugf("GroupUserRelationController: groups_id: %d, user_name: %s\n", groupID, searchUserName)

	if groupID == 0 {
		err := fmt.Errorf("get group available members failed, err: group is invalid, groupID: %d", groupID)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	users, count, err := ctl.svc.FindAvailableMembers(uint64(groupID), searchUserName, page, pageSize)

	if err != nil {
		log.Errorf("get group(%d) members list failed, user_name: %s err: %s", groupID, searchUserName, err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"count": count, "data": users})
}

// DeleteMembers pare group user relations deleting operations.
func (ctl *GroupController) DeleteMembers(ctx echo.Context) error {

	var gur []notice.GroupUserRelation

	if err := ctx.Bind(&gur); err != nil {
		log.Error("delete group user relations get body failed, err: ", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	if len(gur) == 0 {
		return ctx.String(http.StatusBadRequest, "delete members list lenght is 0")
	}

	for _, one := range gur {
		if one.GroupID == 0 || one.UserID == 0 {
			return ctx.String(http.StatusBadRequest, "must provide group_id and user_id")
		}
	}

	if err := ctl.svc.DeleteMembers(gur); err != nil {
		log.Error("delete group user relations failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, gur)
}

// AddMembers parse group user reations creating operation.
func (ctl *GroupController) AddMembers(ctx echo.Context) error {

	log.Debug("start AddMembers")

	var gur []notice.GroupUserRelation

	if err := ctx.Bind(&gur); err != nil {
		log.Error("create group user relations get body failed, err: ", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	log.Debugf("AddMembers: %#v\n", gur)

	if len(gur) == 0 {
		return ctx.String(http.StatusBadRequest, "add members list lenght is 0")
	}

	for _, one := range gur {
		if one.GroupID == 0 || one.UserID == 0 {
			return ctx.String(http.StatusBadRequest, "must provide group_id and user_id")
		}
	}

	if err := ctl.svc.AddMembers(gur); err != nil {
		log.Error("create group user relations failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, gur)
}
