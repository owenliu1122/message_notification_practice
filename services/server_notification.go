package services

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"message_notification_practice/model"
	"message_notification_practice/mq"
)

func NewNotificationService(db *gorm.DB, mq *mq.BaseMq) *NotificationService {
	return &NotificationService{
		db: db,
		mq: mq,
	}
}

type NotificationService struct {
	db *gorm.DB
	mq *mq.BaseMq
}

func (u *NotificationService) Create(notify *model.NotificationRecord) error {
	var err error
	var jsonBytes []byte

	jsonBytes, err = json.Marshal(notify)
	if err != nil {
		return err
	}

	err = u.mq.Send("", "", jsonBytes)
	if err != nil {
		return err
	}

	return u.db.Create(notify).Error
}

func (u *NotificationService) Update(user *model.NotificationRecord, fields map[string]interface{}) error {
	panic("not implemented")
}

func (u *NotificationService) Find(id uint) (*model.NotificationRecord, error) {
	panic("not implemented")
}

func (u *NotificationService) FindByName(name string) (*model.NotificationRecord, error) {
	panic("not implemented")
}

func (u *NotificationService) Delete(user *model.NotificationRecord) (*model.NotificationRecord, error) {
	panic("not implemented")
}
