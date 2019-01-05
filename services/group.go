package services

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/fpay/foundation-go/database"

	"github.com/fpay/foundation-go/log"
	"github.com/owenliu1122/notice"
)

// NewGroupService returns group record operation service.
func NewGroupService(logger *log.Logger, db *database.DB, cache notice.Cache) *GroupService {
	return &GroupService{
		logger: logger,
		db:     db,
		cache:  cache,
	}
}

// GroupService is a group record operation service.
type GroupService struct {
	logger *log.Logger
	db     *database.DB
	cache  notice.Cache
}

// Create a group record.
func (svc *GroupService) Create(group *notice.Group) error {
	return svc.db.Create(group).Error
}

// Update a group record.
func (svc *GroupService) Update(group *notice.Group) error {
	return svc.db.Model(group).Updates(*group).Error
}

// List group list by name , page, page size.
func (svc *GroupService) List(name string, page, pageSize int) (users []notice.Group, count int, err error) {
	group := notice.Group{}
	err = svc.db.Model(&group).
		Where("name like ?", name+"%").
		Count(&count).
		Order("id").
		Offset((page - 1) * pageSize).
		Limit(page * pageSize).
		Find(&users).
		Error

	return
}

// Find a group record by id.
func (svc *GroupService) Find(id uint) (*notice.Group, error) {

	var group notice.Group

	err := svc.db.Find(&group).Error
	//err := u.db.Raw("select * from groups").Scan(&groups).Error

	return &group, err
}

// FindByName a group record by name.
func (svc *GroupService) FindByName(name string) (*notice.Group, error) {
	panic("not implemented")
}

// Delete a group record.
func (svc *GroupService) Delete(group *notice.Group) (*notice.Group, error) {
	return group, svc.db.Delete(group).Error
}

// AddMembers a group user relation record.
func (svc *GroupService) AddMembers(gur []notice.GroupUserRelation) error {
	var err error
	tx := svc.db.Begin()

	for _, one := range gur {
		svc.logger.Debugf("INSERT: %v\n", one)

		cacheKey := getGroupUsersCacheKey(one.GroupID)

		err = svc.db.Create(&one).Error
		if err != nil {
			svc.logger.Error("group user relations create failed, err: ", err)
			svc.cache.Delete(cacheKey)
			break
		}
		user := notice.User{ID: one.UserID}

		if err = svc.db.Find(&user).Error; err != nil {
			svc.logger.Error("group user find user failed, err: ", err)
			svc.cache.Delete(cacheKey)
			break
		}

		if err = svc.cache.SAdd(cacheKey, &user); err != nil {
			svc.logger.Errorf("group user SAdd user failed, cacheKey: %s, err: %s", cacheKey, err)
			svc.cache.Delete(cacheKey)
		}
	}

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return err
}

// FindMembers returns all users that belong with the group id.
func (svc *GroupService) FindMembers(id uint64) ([]notice.User, error) {
	var users []notice.User
	var err error
	//var count int
	// TODO: 未实现分页
	gur := notice.GroupUserRelation{GroupID: id}

	cacheKey := getGroupUsersCacheKey(gur.GroupID)
	if svc.cache.IsExist(cacheKey) {
		if err = svc.cache.SMembers(cacheKey, &users); err != nil {
			svc.cache.Delete(cacheKey)
			svc.logger.Errorf("FindMembers cache smembers failed, group_id: %d, err: %s", id, err)
		} else {
			svc.logger.Debugf("Find Members from cahce, group_id: %d\n", gur.GroupID)
			return users, nil
		}
	}

	err = svc.db.Where("id in (?)", svc.db.Model(&gur).Where(gur).Select("user_id").QueryExpr()).Find(&users).Error
	//err := u.db.Raw("select * from users where id in (select user_id from group_user_relations where group_id = ?)", id).Scan(&users).Error
	if err != nil {
		svc.logger.Errorf("get group(%d) members failed, err: %s\n", id, err)
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.Wrapf(errors.New("group no members"), "group_id: %d", gur.GroupID)
	}

	s := make([]interface{}, len(users))
	for i, v := range users {
		s[i] = v
	}

	if err = svc.cache.SAdd(cacheKey, s...); err != nil {
		svc.logger.Errorf("FindMembers cache SAdd failed, group_id: %d, err: %s", id, err)
	}

	svc.logger.Debugf("Find Members from mysql, group_id: %d\n", gur.GroupID)

	return users, err
}

// FindAvailableMembers will list all users that can be add to current group id.
func (svc *GroupService) FindAvailableMembers(id uint64, uname string, page, pageSize int) ([]notice.User, int, error) {
	var count int
	var users []notice.User
	gur := notice.GroupUserRelation{GroupID: id}

	err := svc.db.Model(&notice.User{}).
		Where("id not in (?) and name like ?", svc.db.Model(&gur).Where(gur).Select("user_id").QueryExpr(), uname+"%").
		Count(&count).
		Offset((page - 1) * pageSize).
		Limit(page * pageSize).
		Find(&users).
		Error

	if err != nil {
		svc.logger.Errorf("get group(%d) members failed, user_name: %serr: %s\n", id, uname, err)
	}

	return users, count, err
}

// DeleteMembers pare group user relations deleting operations.
func (svc *GroupService) DeleteMembers(gur []notice.GroupUserRelation) error {

	tx := svc.db.Begin()

	for _, one := range gur {
		svc.logger.Debugf("DELETE: %v\n", one)

		if err := svc.db.Where(one).Delete(notice.GroupUserRelation{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		svc.cache.Delete(getGroupUsersCacheKey(one.GroupID))

	}

	tx.Commit()

	return nil
}

func getGroupUsersCacheKey(groupID uint64) string {
	return fmt.Sprintf("group_user_relations_group_users_%d", groupID)
}
