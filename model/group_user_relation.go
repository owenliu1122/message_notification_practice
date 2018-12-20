package model

import (
	"time"
)

type GroupUserRelation struct {
	ID        uint64
	GroupID   uint64
	UserID    uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type GroupUserRelationService interface {
	Create(user *GroupUserRelation) error
	Update(user *GroupUserRelation, fields map[string]interface{}) error
	Find(id uint64) (*GroupUserRelation, error)
	FindByName(name string) (*GroupUserRelation, error)
	Delete(user *GroupUserRelation) error
}
