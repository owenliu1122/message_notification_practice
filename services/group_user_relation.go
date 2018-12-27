package services

import (
	"github.com/jinzhu/gorm"
	log "gopkg.in/cihub/seelog.v2"
	"message_notification_practice"
)

func NewGroupUserRelationService(db *gorm.DB) *GroupUserRelationService {
	return &GroupUserRelationService{db: db}
}

type GroupUserRelationService struct {
	db *gorm.DB
}

func (u *GroupUserRelationService) Create(gur []root.GroupUserRelation) error {
	tx := u.db.Begin()

	for _, one := range gur {
		log.Debugf("INSERT: %v\n", one)
		if err := u.db.Create(&one).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (u *GroupUserRelationService) Update(gur *root.GroupUserRelation, fields map[string]interface{}) error {
	panic("not implemented")
}

func (u *GroupUserRelationService) FindMembers(id uint64) ([]root.User, error) {
	var users []root.User
	gur := root.GroupUserRelation{GroupID: id}

	err := u.db.Where("id in (?)", u.db.Model(&gur).Where(gur).Select("user_id").QueryExpr()).Find(&users).Error
	//err := u.db.Raw("select * from users where id in (select user_id from group_user_relations where group_id = ?)", id).Scan(&users).Error
	if err != nil {
		log.Errorf("get group(%d) members failed, err: %s\n", id, err)
	}

	return users, err
}

func (u *GroupUserRelationService) FindAvailableMembers(id uint64, uname string) ([]root.User, error) {
	var users []root.User
	gur := root.GroupUserRelation{GroupID: id}

	//expr:=u.db.Where ("id not in (?) and name like ?",
	//	u.db.Model(&gur).Where(gur).Select("user_id").QueryExpr(),
	//	"%"+uname).QueryExpr()
	//log.Debug("FindAvailableMembers: expr: %#v\n", expr)

	err := u.db.Where("id not in (?) and name like ?",
		u.db.Model(&gur).Where(gur).Select("user_id").QueryExpr(),
		uname+"%").Find(&users).Error
	//err := u.db.Where ("name LIKE ?", "%sh%").Find(&users).Error

	if err != nil {
		log.Errorf("get group(%d) members failed, user_name: %serr: %s\n", id, uname, err)
	}

	return users, err
}

func (u *GroupUserRelationService) FindByName(name string) (*root.GroupUserRelation, error) {
	panic("not implemented")
}

func (u *GroupUserRelationService) Delete(gur []root.GroupUserRelation) error {

	tx := u.db.Begin()

	for _, one := range gur {
		log.Debugf("DELETE: %v\n", one)
		if err := u.db.Where(one).Delete(root.GroupUserRelation{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}
