package model

import "github.com/jinzhu/gorm"

// 用户 模型
type User struct {
	gorm.Model

	User string
	Email string
	Tel string
	Password string
	PasswordSalt string
	Token string

	RoleID uint
	Role Role

	Cart Cart
}
