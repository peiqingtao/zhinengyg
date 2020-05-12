package model

import "github.com/jinzhu/gorm"

type Address struct {
	gorm.Model

	UserID uint
	Tag string // 家，公司，别墅，独栋

	ProvinceCode string `gorm:"type:char(6)"`
	Province string // 省
	CityCode string `gorm:"type:char(6)"`
	City string // 市
	CountyCode string `gorm:"type:char(6)" form:"areaCode" json:"areaCode"`
	County string // 区

	Addr string `form:"addressDetail" json:"addressDetail"`// 详细地址，U8产业园U6-500

	IsDefault bool

	Name string // 收货人名字
	Tel string // 收货人电话

	PostCode string `form:"postalCode" json:"postalCode"` // 邮政编码
}
