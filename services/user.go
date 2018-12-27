package services

import (
	"github.com/jinzhu/gorm"
	"message_notification_practice"
)

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

type UserService struct {
	db *gorm.DB
}

func (u *UserService) Create(user *root.User) error {
	return u.db.Create(user).Error
}

func (u *UserService) Update(user *root.User, fields map[string]interface{}) error {
	return u.db.Model(user).Updates(*user).Error
}

func (u *UserService) Find(id uint) ([]root.User, error) {
	var users []root.User

	err := u.db.Find(&users).Error
	//err := u.db.Raw("select * from groups").Scan(&groups).Error

	return users, err
}

func (u *UserService) FindByName(name string) (*root.User, error) {
	panic("not implemented")
}

func (u *UserService) Delete(user *root.User) (*root.User, error) {
	return user, u.db.Delete(user).Error
}
