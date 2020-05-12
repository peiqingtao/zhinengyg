package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Order struct {
	gorm.Model

	Sn string
	AddressID uint
	UserID uint
	OrderStatusID uint
	PaymentID uint
	PaymentStatusID uint
	PaymentSn string
	ShippingID uint
	ShippingStatusID uint
	ShippingSn string
	ShippingAmount int
	Amount int
	TaxID uint
	TaxAmount int
	ProductAmount int
	OrderTime time.Time // 订单的确认时间
	Note string
}
