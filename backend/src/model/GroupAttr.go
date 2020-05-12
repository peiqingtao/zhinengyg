package model

import "github.com/jinzhu/gorm"

// 分组差异属性 模型
type GroupAttr struct {
	gorm.Model

	GroupID uint
	AttrID uint
}
