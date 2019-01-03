package controllers

import (
	"net/http"

	"github.com/owenliu1122/notice"
	"github.com/owenliu1122/notice/services"

	"github.com/labstack/echo"
	log "gopkg.in/cihub/seelog.v2"
)

// NewUserController returns an user table operation controller.
func NewUserController(us *services.UserService) *UserController {
	return &UserController{svc: us}
}

// UserController is an user table operation controller.
type UserController struct {
	svc *services.UserService
}

// List will return all users in the users table.
func (ctl *UserController) List(ctx echo.Context) error {

	res, err := ctl.svc.Find(0)

	if err != nil {
		log.Error("get users list failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, res)
}

// Create will insert a new user record.
func (ctl *UserController) Create(ctx echo.Context) error {

	var user notice.User
	if err := ctx.Bind(&user); err != nil {
		log.Error("add user get body failed, err: ", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	log.Infof("UserController Bind -> user: %v\n", user)

	if user.Name == "" {
		err := log.Error("create user failed, err: no user name.")
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	if user.Email == "" && user.Phone == "" && user.Wechat == "" {
		err := log.Error("create user failed, did not fill in any communication method")
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err := ctl.svc.Create(&user)

	if err != nil {
		log.Error("create user failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, user)
}

// Update will update an user record.
func (ctl *UserController) Update(ctx echo.Context) error {

	var user notice.User
	if err := ctx.Bind(&user); err != nil {
		log.Error("update group get body failed, err: ", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	log.Infof("GroupController Update -> user: %#v\n", user)

	if user.ID == 0 {
		err := log.Error("update user failed, err: no user id.")
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	if user.Name == "" && user.Email == "" && user.Phone == "" && user.Wechat == "" {
		err := log.Error("update user failed, did not fill in any modification information")
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	err := ctl.svc.Update(&user)

	if err != nil {
		log.Error("update user failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, user)
}

// Delete will delete an user record.
func (ctl *UserController) Delete(ctx echo.Context) error {

	var user notice.User
	if err := ctx.Bind(&user); err != nil {
		log.Error("delete user get body failed, err: ", err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	log.Infof("GroupController Delete -> user: %#v\n", user)

	if user.ID == 0 {
		err := log.Error("deleted user failed, err: no user id.")
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	_, err := ctl.svc.Delete(&user)

	if err != nil {
		log.Error("delete user failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, user)
}
