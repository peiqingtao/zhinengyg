package model

import "github.com/jinzhu/gorm"

// 产品属性 模型
type ProductAttr struct {
	gorm.Model

	ProductID uint
	AttrID uint
	Value string
	SortOrder int

}
