package model

import "github.com/jinzhu/gorm"

// 属性类型 模型
type AttrType struct {
	gorm.Model

	Name string
	SortOrder int

	// 关联定义
	AttrGroups []AttrGroup
}
