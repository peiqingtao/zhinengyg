package model

import "github.com/jinzhu/gorm"

// 分组 模型
type Group struct {
	gorm.Model

	Counter int
	Name string
	SortOrder int
	AttrTypeID uint

	CheckedProductID []uint `gorm:"-"`
	CheckedAttrID []uint `gorm:"-"`

	Products []Product
	Attrs []Attr `gorm:"many2many:group_attrs"`

}
