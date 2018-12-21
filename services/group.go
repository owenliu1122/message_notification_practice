package services

import (
	"github.com/jinzhu/gorm"
	"message_notification_practice/model"
)

func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{db: db}
}

type GroupService struct {
	db *gorm.DB
}

func (u *GroupService) Create(user *model.Group) error {
	panic("not implemented")
}

func (u *GroupService) Update(user *model.Group, fields map[string]interface{}) error {
	panic("not implemented")
}

func (u *GroupService) Find(id uint) ([]model.Group, error) {

	var groups []model.Group

	err := u.db.Raw("select * from groups").Scan(&groups).Error

	return groups, err
}

func (u *GroupService) FindByName(name string) (*model.Group, error) {
	panic("not implemented")
}

func (u *GroupService) Delete(user *model.Group) (*model.Group, error) {
	panic("not implemented")
}
