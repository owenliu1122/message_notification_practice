package model

import "time"

type Group struct {
	ID        uint64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type GroupService interface {
	Create(group *Group) error
	Update(user *Group, fields map[string]interface{}) error
	Find(id uint64) (*Group, error)
	FindByName(name string) (*Group, error)
	Delete(user *Group) error
}
