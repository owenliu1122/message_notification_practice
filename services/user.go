package services

import "message_notification_practice/model"

func NewUserService() *UserService {
	return &UserService{}
}

type UserService struct{}

func (u *UserService) Create(user *model.User) error {
	panic("not implemented")
}

func (u *UserService) Update(user *model.User, fields map[string]interface{}) error {
	panic("not implemented")
}

func (u *UserService) Find(id uint) (*model.User, error) {
	panic("not implemented")
}

func (u *UserService) FindByName(name string) (*model.User, error) {
	panic("not implemented")
}

func (u *UserService) Delete(user *model.User) (*model.User, error) {
	panic("not implemented")
}
