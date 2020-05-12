package model

import "github.com/jinzhu/gorm"

type Quantity struct {
	gorm.Model
	ProductID uint
	Number int
	StoreID uint
}
