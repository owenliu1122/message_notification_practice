package notice

import (
	"time"
)

// Group is an user group infomation.
type Group struct {
	ID        uint64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Name      string     `gorm:"column:name;not null;unique" json:"name"`
	CreatedAt time.Time  `gorm:"column:created_at;not null;DEFAULT:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null;DEFAULT:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// GroupServiceInterface is an user group method interface.
type GroupServiceInterface interface {
	Create(group *Group) error
	Update(group *Group) error
	List(name string, page, pageSize int) ([]Group, int, error)
	Find(id uint) (*Group, error)
	FindByName(name string) (*Group, error)
	Delete(group *Group) (*Group, error)
	AddMembers(gur []GroupUserRelation) error
	FindMembers(id uint64) ([]User, error)
	FindAvailableMembers(id uint64, uname string, page, pageSize int) ([]User, int, error)
	DeleteMembers(gur []GroupUserRelation) error
}
