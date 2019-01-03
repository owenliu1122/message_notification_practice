package services

import (
	"github.com/jinzhu/gorm"
	"github.com/owenliu1122/notice"
	log "gopkg.in/cihub/seelog.v2"
)

// NewUserService returns user record operation service.
func NewUserService(db *gorm.DB, cache notice.Cache) *UserService {
	return &UserService{
		db:    db,
		cache: cache,
	}
}

// UserService is a user record operation service.
type UserService struct {
	db    *gorm.DB
	cache notice.Cache
}

// Create a user record.
func (u *UserService) Create(user *notice.User) error {
	return u.db.Create(user).Error
}

// Update a user record.
func (u *UserService) Update(user *notice.User) error {
	if err := u.db.Model(user).Updates(*user).Error; err != nil {
		return err
	}

	u.deleteGroupCaches(user)

	return nil
}

// Find a user record by id.
func (u *UserService) Find(id uint) ([]notice.User, error) {
	var users []notice.User

	err := u.db.Find(&users).Error
	//err := u.db.Raw("select * from groups").Scan(&groups).Error

	return users, err
}

// FindByName a user record by name.
func (u *UserService) FindByName(name string) (*notice.User, error) {
	panic("not implemented")
}

// Delete a user record.
func (u *UserService) Delete(user *notice.User) (*notice.User, error) {
	if err := u.db.Delete(user).Error; err != nil {
		return nil, err
	}

	u.deleteGroupCaches(user)

	return user, nil
}

func (u *UserService) deleteGroupCaches(user *notice.User) error {
	var gurs []notice.GroupUserRelation

	if err := u.db.Where("user_id = ?", user.ID).Select("group_id").Find(&gurs).Error; err != nil {
		log.Errorf("user update find belone group id failed, err: ", err)
		return err
	}

	// 这个用户不属于任何 group
	if len(gurs) == 0 {
		return nil
	}

	var gurCacheKeys []string

	for _, one := range gurs {
		gurCacheKeys = append(gurCacheKeys, getGroupUsersCacheKey(one.GroupID))
	}

	if err := u.cache.Delete(gurCacheKeys...); err != nil {
		log.Error("delete related group cache failed, err: ", err)
		return err
	}

	return nil
}
