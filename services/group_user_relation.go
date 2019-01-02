package services

import (
	"fmt"
	"github.com/owenliu1122/notice"

	"github.com/owenliu1122/notice/redis"

	"github.com/jinzhu/gorm"
	log "gopkg.in/cihub/seelog.v2"
)

// NewGroupUserRelationService return group user relation operation service.
func NewGroupUserRelationService(db *gorm.DB, cache redis.Cache) *GroupUserRelationService {
	return &GroupUserRelationService{
		db:    db,
		cache: cache,
	}
}

// GroupUserRelationService is a group user relation service.
type GroupUserRelationService struct {
	db    *gorm.DB
	cache redis.Cache
}

// Create a group user relation record.
func (svc *GroupUserRelationService) Create(gur []notice.GroupUserRelation) error {
	var err error
	tx := svc.db.Begin()

	for _, one := range gur {
		log.Debugf("INSERT: %v\n", one)

		cacheKey := getGroupUsersCacheKey(one.GroupID)

		err = svc.db.Create(&one).Error
		if err != nil {
			log.Error("group user relations create failed, err: ", err)
			svc.cache.Delete(cacheKey)
			break
		}
		user := notice.User{ID: one.UserID}

		if err = svc.db.Find(&user).Error; err != nil {
			log.Error("group user find user failed, err: ", err)
			svc.cache.Delete(cacheKey)
			break
		}

		if err = svc.cache.SAdd(cacheKey, &user); err != nil {
			log.Errorf("group user SAdd user failed, cacheKey: %s, err: %s", cacheKey, err)
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

// Update a group user relation record.
func (svc *GroupUserRelationService) Update(gur *notice.GroupUserRelation, fields map[string]interface{}) error {
	svc.cache.Delete(getGroupUsersCacheKey(gur.GroupID))
	panic("not implemented")
}

// FindMembers returns all users that belong with the group id.
func (svc *GroupUserRelationService) FindMembers(id uint64) ([]notice.User, error) {
	var users []notice.User
	var err error

	gur := notice.GroupUserRelation{GroupID: id}

	cacheKey := getGroupUsersCacheKey(gur.GroupID)
	if svc.cache.IsExist(cacheKey) {
		if err = svc.cache.SMembers(cacheKey, &users); err != nil {
			svc.cache.Delete(cacheKey)
			log.Errorf("FindMembers cache smembers failed, group_id: %d, err: %s", id, err)
		} else {
			log.Debugf("Find Members from cahce, group_id: %d\n", gur.GroupID)
			return users, nil
		}
	}

	err = svc.db.Where("id in (?)", svc.db.Model(&gur).Where(gur).Select("user_id").QueryExpr()).Find(&users).Error
	//err := u.db.Raw("select * from users where id in (select user_id from group_user_relations where group_id = ?)", id).Scan(&users).Error
	if err != nil {
		log.Errorf("get group(%d) members failed, err: %s\n", id, err)
	}

	s := make([]interface{}, len(users))
	for i, v := range users {
		s[i] = v
	}

	if err = svc.cache.SAdd(cacheKey, s...); err != nil {
		log.Errorf("FindMembers cache SAdd failed, group_id: %d, err: %s", id, err)
	}

	log.Debugf("Find Members from mysql, group_id: %d\n", gur.GroupID)

	return users, err
}

// FindAvailableMembers will list all users that can be add to current group id.
func (svc *GroupUserRelationService) FindAvailableMembers(id uint64, uname string) ([]notice.User, error) {
	var users []notice.User
	gur := notice.GroupUserRelation{GroupID: id}

	//expr:=u.db.Where ("id not in (?) and name like ?",
	//	u.db.Model(&gur).Where(gur).Select("user_id").QueryExpr(),
	//	"%"+uname).QueryExpr()
	//log.Debug("FindAvailableMembers: expr: %#v\n", expr)

	err := svc.db.Where("id not in (?) and name like ?",
		svc.db.Model(&gur).Where(gur).Select("user_id").QueryExpr(),
		uname+"%").Find(&users).Error
	//err := u.db.Where ("name LIKE ?", "%sh%").Find(&users).Error

	if err != nil {
		log.Errorf("get group(%d) members failed, user_name: %serr: %s\n", id, uname, err)
	}

	return users, err
}

// Delete pare group user relations deleting operations.
func (svc *GroupUserRelationService) Delete(gur []notice.GroupUserRelation) error {

	tx := svc.db.Begin()

	for _, one := range gur {
		log.Debugf("DELETE: %v\n", one)

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
