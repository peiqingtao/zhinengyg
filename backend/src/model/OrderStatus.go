package model

import "github.com/jinzhu/gorm"

type OrderStatus struct {
	gorm.Model
	Title string
}
