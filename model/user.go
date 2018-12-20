package model

import "time"

type User struct {
	ID        uint64 `gorm:"column:id" json:"id"`
	Name      string
	Phone     string
	Email     string
	Wechat    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type UserService interface {
	Create(user *User) error
	Update(user *User, fields map[string]interface{}) error
	Find(id uint64) (*User, error)
	FindByName(name string) (*User, error)
	Delete(user *User) error
}
