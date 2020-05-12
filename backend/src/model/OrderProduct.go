package model

import "github.com/jinzhu/gorm"

type OrderProduct struct {
	gorm.Model

	OrderID uint
	ProductID uint
	BuyQuantity int
	BuyPrice int
}
