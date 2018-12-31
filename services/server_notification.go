package services

import (
	"encoding/json"
	"message_notification_practice"
	"message_notification_practice/mq"
	"message_notification_practice/pb"
	"strings"

	"github.com/jinzhu/gorm"
)

// NewNotificationService returns a notification record operation service.
func NewNotificationService(db *gorm.DB, mq *mq.BaseMq) *NotificationService {
	return &NotificationService{
		db: db,
		mq: mq,
	}
}

// NotificationService is a notification record operation service.
type NotificationService struct {
	db *gorm.DB
	mq *mq.BaseMq
}

// Create a notification record.
func (u *NotificationService) Create(pbReq *pb.MsgNotificationRequest) error {
	var err error
	var jsonBytes []byte

	typesStr := make([]string, len(pbReq.NoticeType))
	for _, noticeType := range pbReq.NoticeType {
		typesStr = append(typesStr, noticeType.String())
	}

	if err = u.db.Create(&notice.NotificationRecord{
		GroupID:      pbReq.Group,
		NoticeType:   strings.Join(typesStr, ","),
		Notification: pbReq.Content,
	}).Error; err != nil {
		return err
	}

	jsonBytes, err = json.Marshal(pbReq)
	if err != nil {
		return err
	}

	err = u.mq.Send("", "", jsonBytes)
	if err != nil {
		return err
	}

	return err
}

// Update notification records.
func (u *NotificationService) Update(user *notice.NotificationRecord, fields map[string]interface{}) error {
	panic("not implemented")
}

// Find a notification record.
func (u *NotificationService) Find(id uint) (*notice.NotificationRecord, error) {
	panic("not implemented")
}

// FindByName a notification record.
func (u *NotificationService) FindByName(name string) (*notice.NotificationRecord, error) {
	panic("not implemented")
}

// Delete a notification record.
func (u *NotificationService) Delete(user *notice.NotificationRecord) (*notice.NotificationRecord, error) {
	panic("not implemented")
}
