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

func (u *GroupService) Create(group *model.Group) error {
	return u.db.Create(group).Error
}

func (u *GroupService) Update(group *model.Group, fields map[string]interface{}) error {
	return u.db.Model(group).Updates(*group).Error
}

func (u *GroupService) Find(id uint) ([]model.Group, error) {

	var groups []model.Group

	err := u.db.Find(&groups).Error
	//err := u.db.Raw("select * from groups").Scan(&groups).Error

	return groups, err
}

func (u *GroupService) FindByName(name string) (*model.Group, error) {
	panic("not implemented")
}

func (u *GroupService) Delete(group *model.Group) (*model.Group, error) {
	return group, u.db.Delete(group).Error
}
