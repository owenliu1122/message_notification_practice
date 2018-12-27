package root

import "time"

type User struct {
	ID        uint64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Name      string     `gorm:"column:name;not null;unique" json:"name"`
	Phone     string     `gorm:"column:phone" json:"phone"`
	Email     string     `gorm:"column:email" json:"email"`
	Wechat    string     `gorm:"column:wechat" json:"wechat"`
	CreatedAt time.Time  `gorm:"column:created_at;not null;DEFAULT:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;not null;DEFAULT:CURRENT_TIMESTAMP"json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

type UserService interface {
	Create(user *User) error
	Update(user *User, fields map[string]interface{}) error
	Find(id uint64) (*User, error)
	FindByName(name string) (*User, error)
	Delete(user *User) error
}
