package services

import "message_notification_practice/model"

func NewGroupUserRelationService() *GroupUserRelationService {
	return &GroupUserRelationService{}
}

type GroupUserRelationService struct{}

func (u *GroupUserRelationService) Create(user *model.GroupUserRelation) error {
	panic("not implemented")
}

func (u *GroupUserRelationService) Update(user *model.GroupUserRelation, fields map[string]interface{}) error {
	panic("not implemented")
}

func (u *GroupUserRelationService) Find(id uint) (*model.GroupUserRelation, error) {
	panic("not implemented")
}

func (u *GroupUserRelationService) FindByName(name string) (*model.GroupUserRelation, error) {
	panic("not implemented")
}

func (u *GroupUserRelationService) Delete(user *model.GroupUserRelation) (*model.GroupUserRelation, error) {
	panic("not implemented")
}
