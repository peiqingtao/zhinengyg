package model

import "github.com/jinzhu/gorm"

// 角色授权 模型
type RolePrivilege struct {
	gorm.Model

	RoleID uint
	PrivilegeID uint

}
