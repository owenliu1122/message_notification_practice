package services

import (
	"github.com/jinzhu/gorm"
)

// NewGroupService returns group record operation service.
func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{db: db}
}

// GroupService is a group record operation service.
type GroupService struct {
	db *gorm.DB
}

// Create a group record.
func (u *GroupService) Create(group *notice.Group) error {
	return u.db.Create(group).Error
}

// Update a group record.
func (u *GroupService) Update(group *notice.Group) error {
	return u.db.Model(group).Updates(*group).Error
}

// Find a group record by id.
func (u *GroupService) Find(id uint) ([]notice.Group, error) {

	var groups []notice.Group

	err := u.db.Find(&groups).Error
	//err := u.db.Raw("select * from groups").Scan(&groups).Error

	return groups, err
}

// FindByName a group record by name.
func (u *GroupService) FindByName(name string) (*notice.Group, error) {
	panic("not implemented")
}

// Delete a group record.
func (u *GroupService) Delete(group *notice.Group) (*notice.Group, error) {
	return group, u.db.Delete(group).Error
}
