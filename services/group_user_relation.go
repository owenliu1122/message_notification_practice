package services

import (
	"fmt"
	"github.com/jinzhu/gorm"
	log "gopkg.in/cihub/seelog.v2"
	"message_notification_practice"
	"message_notification_practice/redis"
)

func NewGroupUserRelationService(db *gorm.DB, cache redis.Cache) *GroupUserRelationService {
	return &GroupUserRelationService{
		db:    db,
		cache: cache,
	}
}

type GroupUserRelationService struct {
	db    *gorm.DB
	cache redis.Cache
}

func (svc *GroupUserRelationService) Create(gur []root.GroupUserRelation) error {
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
		user := root.User{ID: 1}

		if err = svc.db.Find(&user).Error; err != nil {
			log.Error("group user find user failed, err: ", err)
			svc.cache.Delete(cacheKey)
			break
		}

		if err := svc.cache.SAdd(cacheKey, &user); err != nil {
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

func (svc *GroupUserRelationService) Update(gur *root.GroupUserRelation, fields map[string]interface{}) error {
	panic("not implemented")
	svc.cache.Delete(getGroupUsersCacheKey(gur.GroupID))
	return nil
}

func (svc *GroupUserRelationService) FindMembers(id uint64) ([]root.User, error) {
	var users []root.User
	gur := root.GroupUserRelation{GroupID: id}

	cacheKey := getGroupUsersCacheKey(gur.GroupID)
	if svc.cache.IsExist(cacheKey) {
		if err := svc.cache.SMembers(cacheKey, &users); err != nil {
			svc.cache.Delete(cacheKey)
			log.Errorf("FindMembers cache smembers failed, group_id: %d, err: %s", id, err)
		} else {
			log.Debugf("Find Members from cahce, group_id: %d\n", gur.GroupID)
			return users, nil
		}
	}

	err := svc.db.Where("id in (?)", svc.db.Model(&gur).Where(gur).Select("user_id").QueryExpr()).Find(&users).Error
	//err := u.db.Raw("select * from users where id in (select user_id from group_user_relations where group_id = ?)", id).Scan(&users).Error
	if err != nil {
		log.Errorf("get group(%d) members failed, err: %s\n", id, err)
	}

	s := make([]interface{}, len(users))
	for i, v := range users {
		s[i] = v
	}

	if err := svc.cache.SAdd(cacheKey, s...); err != nil {
		log.Errorf("FindMembers cache SAdd failed, group_id: %d, err: %s", id, err)
	}

	log.Debugf("Find Members from mysql, group_id: %d\n", gur.GroupID)

	return users, err
}

func (svc *GroupUserRelationService) FindAvailableMembers(id uint64, uname string) ([]root.User, error) {
	var users []root.User
	gur := root.GroupUserRelation{GroupID: id}

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

func (svc *GroupUserRelationService) FindByName(name string) (*root.GroupUserRelation, error) {
	panic("not implemented")
}

func (svc *GroupUserRelationService) Delete(gur []root.GroupUserRelation) error {

	tx := svc.db.Begin()

	for _, one := range gur {
		log.Debugf("DELETE: %v\n", one)

		if err := svc.db.Where(one).Delete(root.GroupUserRelation{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		svc.cache.Delete(getGroupUsersCacheKey(one.GroupID))

	}

	tx.Commit()

	return nil
}

func getGroupUsersCacheKey(groupId uint64) string {
	return fmt.Sprintf("group_user_relations_group_users_%d", groupId)
}
