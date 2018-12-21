package controllers

import (
	"message_notification_practice/services"
)

func NewUserController(us *services.UserService) *UserController {
	return &UserController{us: us}
}

type UserController struct {
	us *services.UserService
}
