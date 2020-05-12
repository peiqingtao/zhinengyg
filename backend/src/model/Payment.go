package model

import "github.com/jinzhu/gorm"

type Payment struct {
	gorm.Model
	Title string
	Key string
	Intro string
	Status int // 0 1 2 分别表示 禁用，启用，维护 等状态
}
