package services

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"message_notification_practice"
	"message_notification_practice/redis"
)

func NewUserService(db *gorm.DB, cache redis.Cache) *UserService {
	return &UserService{
		db:    db,
		cache: cache,
	}
}

type UserService struct {
	db    *gorm.DB
	cache redis.Cache
}

func (u *UserService) Create(user *root.User) error {
	return u.db.Create(user).Error
}

func (u *UserService) Update(user *root.User, fields map[string]interface{}) error {
	if err := u.db.Model(user).Updates(*user).Error; err != nil {
		return err
	}

	u.deleteGroupCaches(user)

	return nil
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
	if err := u.db.Delete(user).Error; err != nil {
		return nil, err
	}

	u.deleteGroupCaches(user)

	return user, nil
}

func (u *UserService) deleteGroupCaches(user *root.User) error {
	var gurs []root.GroupUserRelation

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
