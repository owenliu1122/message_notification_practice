package services

import (
	"encoding/json"
	"strings"

	"github.com/owenliu1122/notice"
	"github.com/owenliu1122/notice/pb"

	"github.com/jinzhu/gorm"
)

// NewNotificationService returns a notification record operation service.
func NewNotificationService(db *gorm.DB, pc notice.ProducerInterface, exchange, routing string) *NotificationService {
	return &NotificationService{
		db:         db,
		pc:         pc,
		pcExchange: exchange,
		pcRouting:  routing,
	}
}

// NotificationService is a notification record operation service.
type NotificationService struct {
	db         *gorm.DB
	pc         notice.ProducerInterface
	pcExchange string
	pcRouting  string
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

	err = u.pc.Publish(u.pcExchange, u.pcRouting, jsonBytes)
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
