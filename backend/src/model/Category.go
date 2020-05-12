package model

import "github.com/jinzhu/gorm"

// 模型结构体(类)定义
type Category struct {
	// 嵌套结构体
	gorm.Model

	ParentId uint
	Name string
	Logo string
	Description string
	SortOrder int
	MetaTitle string
	MetaKeywords string
	MetaDescription string
}