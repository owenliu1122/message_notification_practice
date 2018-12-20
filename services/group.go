package services

import "message_notification_practice/model"

func NewGroupService() *GroupService {
	return &GroupService{}
}

type GroupService struct{}

func (u *GroupService) Create(user *model.Group) error {
	panic("not implemented")
}

func (u *GroupService) Update(user *model.Group, fields map[string]interface{}) error {
	panic("not implemented")
}

func (u *GroupService) Find(id uint) (*model.Group, error) {
	panic("not implemented")
}

func (u *GroupService) FindByName(name string) (*model.Group, error) {
	panic("not implemented")
}

func (u *GroupService) Delete(user *model.Group) (*model.Group, error) {
	panic("not implemented")
}
