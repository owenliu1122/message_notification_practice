package services

import (
	"github.com/fpay/foundation-go/database"
	"github.com/fpay/foundation-go/log"
	"github.com/owenliu1122/notice"
)

// NewUserService returns user record operation service.
func NewUserService(logger *log.Logger, db *database.DB, cache notice.Cache) *UserService {
	return &UserService{
		logger: logger,
		db:     db,
		cache:  cache,
	}
}

// UserService is a user record operation service.
type UserService struct {
	logger *log.Logger
	db     *database.DB
	cache  notice.Cache
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

	return svc.deleteGroupCaches(user)
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
	var err error
	if err = svc.db.Delete(user).Error; err != nil {
		return err
	}

	if err = svc.deleteGroupCaches(user); err != nil {
		svc.logger.Error("delete group cahces failed, err:", err)
	}

	return err
}

func (svc *UserService) deleteGroupCaches(user *notice.User) error {
	var gurs []notice.GroupUserRelation

	if err := svc.db.Where("user_id = ?", user.ID).Select("group_id").Find(&gurs).Error; err != nil {
		svc.logger.Errorf("user delete group id cahce failed, err: ", err)
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
		svc.logger.Error("delete related group cache failed, err: ", err)
		return err
	}

	return nil
}
