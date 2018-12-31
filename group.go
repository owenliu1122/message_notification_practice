package notice

import "time"

// Group is an user group infomation.
type Group struct {
	ID        uint64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Name      string     `gorm:"column:name;not null;unique" json:"name"`
	CreatedAt time.Time  `gorm:"column:created_at;not null;DEFAULT:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null;DEFAULT:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// GroupService is an user group method interface.
type GroupService interface {
	Create(group *Group) error
	Update(user *Group) error
	Find(id uint64) (*Group, error)
	FindByName(name string) (*Group, error)
	Delete(user *Group) error
}
