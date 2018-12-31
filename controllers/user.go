package controllers

import (
	"message_notification_practice"
	"message_notification_practice/services"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
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

	err := ctl.svc.Update(&user, nil)

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

	_, err := ctl.svc.Delete(&user)

	if err != nil {
		log.Error("delete user failed, err: ", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, user)
}
