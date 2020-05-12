package model

import "github.com/jinzhu/gorm"

type PaymentStatus struct {
	gorm.Model
	Title string
}
