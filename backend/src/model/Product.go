package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Product struct {
	gorm.Model

	Name string `gorm:"index"`
	Price float64 `gorm:"type:decimal(14, 2);index"`
	Upc string `gorm:"unique_index;"`
	Mpn string `gorm:"size:127;"`
	IsSale int `gorm:""`
	SaleTime time.Time `gorm:""`
	IsSubstract int `gorm:""`
	IsShipping int `gorm:""`
	Weight float64 `gorm:""`
	Description string `gorm:""`

	// 解析字段
	AttrValue map[uint]string `gorm:"-"`
	//AttrValue map[uint]struct{
	//	Value string
	//	SortOrder int
	//} `gorm:"-"`
	//AttrValue map[uint]map[string]string `gorm:"-"`

	UploadedImage []string `gorm:"-"`
	UploadedImageSmall []string `gorm:"-"`
	UploadedImageBig []string `gorm:"-"`

	// 外键字段
	CategoryID uint
	AttrTypeID uint
	GroupID uint

	// 关联定义
	// 产品属于分类
	Category Category
	AttrType AttrType
	ProductAttrs []ProductAttr
	Images []Image
	Group Group

	// 额外字段
	ModelInfo string `gorm:"-"`
}


