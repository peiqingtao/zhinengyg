package model

import "github.com/jinzhu/gorm"

// 权限 模型
type Privilege struct {
	gorm.Model

	Name string
	Key string
	SortOrder int
	Description string

}
