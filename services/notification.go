package services

import (
	"github.com/jinzhu/gorm"
	"message_notification_practice/model"
)

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{
		db: db,
	}
}

type NotificationService struct {
	db *gorm.DB
}

func (u *NotificationService) Create(user *model.Notification) error {
	return u.db.Create(user).Error
}

func (u *NotificationService) Update(user *model.Notification, fields map[string]interface{}) error {
	panic("not implemented")
}

func (u *NotificationService) Find(id uint) (*model.Notification, error) {
	panic("not implemented")
}

func (u *NotificationService) FindByName(name string) (*model.Notification, error) {
	panic("not implemented")
}

func (u *NotificationService) Delete(user *model.Notification) (*model.Notification, error) {
	panic("not implemented")
}
