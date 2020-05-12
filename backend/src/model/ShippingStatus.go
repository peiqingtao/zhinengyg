package model

import "github.com/jinzhu/gorm"

type ShippingStatus struct {
	gorm.Model
	Title string
}

