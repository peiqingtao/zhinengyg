package model

import "github.com/jinzhu/gorm"

// 属性 模型
type Cart struct {
	gorm.Model

	UserID uint `gorm:"unique_index"`
	Content string `gorm:"type:text"`
}
