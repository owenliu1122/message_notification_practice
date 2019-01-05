package services

import (
	"context"
	"strings"

	"github.com/fpay/foundation-go"

	"github.com/fpay/foundation-go/database"
	"github.com/fpay/foundation-go/log"

	"github.com/owenliu1122/notice"
	"github.com/owenliu1122/notice/pb"
)

// NewNotificationService returns a notification record operation service.
func NewNotificationService(logger *log.Logger, db *database.DB, pc foundation.JobManager, queue string) *NotificationService {
	return &NotificationService{
		logger: logger,
		db:     db,
		pc:     pc,
		job:    &Job{Q: queue},
	}
}

// NotificationService is a notification record operation service.
type NotificationService struct {
	logger *log.Logger
	db     *database.DB
	pc     foundation.JobManager
	job    *Job
}

// Create a notification record.
func (u *NotificationService) Create(ctx context.Context, pbReq *pb.MsgNotificationRequest) error {
	var err error

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

	u.job.Message = pbReq
	err = u.pc.Dispatch(ctx, u.job)
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
