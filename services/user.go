package services

import (
	"github.com/jinzhu/gorm"
	"github.com/owenliu1122/notice"
	log "github.com/sirupsen/logrus"
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
func (svc *UserService) Create(user *notice.User) error {
	return svc.db.Create(user).Error
}

// Update a user record.
func (svc *UserService) Update(user *notice.User) error {
	if err := svc.db.Model(user).Updates(*user).Error; err != nil {
		return err
	}

	svc.deleteGroupCaches(user)

	return nil
}

// List user list by name , page, page size.
func (svc *UserService) List(name string, page, pageSize int) (users []notice.User, count int, err error) {
	user := notice.User{}
	err = svc.db.Model(&user).
		Where("name like ?", name+"%").
		Count(&count).
		Order("id").
		Offset((page - 1) * pageSize).
		Limit(page * pageSize).
		Find(&users).
		Error

	return
}

// Find a user record by id.
func (svc *UserService) Find(id uint) (*notice.User, error) {
	var user notice.User

	err := svc.db.Find(&user).Error
	//err := u.db.Raw("select * from groups").Scan(&groups).Error

	return &user, err
}

// FindByName a user record by name.
func (svc *UserService) FindByName(name string) (*notice.User, error) {
	panic("not implemented")
}

// Delete a user record.
func (svc *UserService) Delete(user *notice.User) error {
	if err := svc.db.Delete(user).Error; err != nil {
		return err
	}

	if err := svc.deleteGroupCaches(user); err != nil {
		log.Error("delete group cahces failed, err:", err)
	}

	return nil
}

func (svc *UserService) deleteGroupCaches(user *notice.User) error {
	var gurs []notice.GroupUserRelation

	if err := svc.db.Where("user_id = ?", user.ID).Select("group_id").Find(&gurs).Error; err != nil {
		log.Errorf("user delete group id cahce failed, err: ", err)
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

	if err := svc.cache.Delete(gurCacheKeys...); err != nil {
		log.Error("delete related group cache failed, err: ", err)
		return err
	}

	return nil
}
