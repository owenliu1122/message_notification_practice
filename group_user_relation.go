package notice

import (
	"time"
)

// GroupUserRelation is group and user relation description information.
type GroupUserRelation struct {
	ID        uint64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	GroupID   uint64     `gorm:"column:group_id;not null" json:"group_id"`
	UserID    uint64     `gorm:"column:user_id;not null" json:"user_id"`
	CreatedAt time.Time  `gorm:"column:created_at;not null;DEFAULT:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null;DEFAULT:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
