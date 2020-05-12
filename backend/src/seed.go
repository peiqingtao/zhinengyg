package main

import (
	"config"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"model"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// 1 启用配置
	config.InitConfig()
	// 2 初始化GORM
	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return config.App["DB_TABLE_PREFIX"] + defaultTableName
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&loc=%s&parseTime=%s",
		config.App["MYSQL_USER"],
		config.App["MYSQL_PASSWORD"],
		config.App["MYSQL_HOST"],
		config.App["MYSQL_PORT"],
		config.App["MYSQL_DBNAME"],
		config.App["MYSQL_CHARSET"],
		config.App["MYSQL_LOC"],
		config.App["MYSQL_PARSETIME"],
	)
	orm, err := gorm.Open(config.App["DB_DRIVER"], dsn)
	if err != nil {
		log.Println(err)
		return
	}
	// 3 seed
	orm.Create(&model.Quantity {
		ProductID: 1, Number: 14,
	})
	orm.Create(&model.Quantity {
		ProductID: 2, Number: 8,
	})
	log.Println("table <quantity> is seeded!")
	//// payment
	//orm.Create(&model.Payment{
	//	Title: "微信支付", Key: "wechat-pay", Intro: "基于微信提供的支付系统", Status: 1,
	//})
	//orm.Create(&model.Payment{
	//	Title: "支付宝支付", Key: "alipay", Intro: "基于支付宝提供的支付系统", Status: 1,
	//})
	//orm.Create(&model.Payment{
	//	Title: "银联", Key: "yinlian-pay", Intro: "基于银联提供的支付系统", Status: 1,
	//})
	//log.Println("table <payment> is seeded!")
	//
	//
	//orm.Create(&model.PaymentStatus{
	//	Title: "支付错误",
	//})
	//orm.Create(&model.PaymentStatus{
	//	Title: "未支付",
	//})
	//orm.Create(&model.PaymentStatus{
	//	Title: "已支付",
	//})
	//log.Println("table <payment_status> is seeded!")
	//
	//
	//// shipping
	//orm.Create(&model.Shipping{
	//	Title: "菜鸟驿站", Key: "cainiao", Intro: "基于微信提供的支付系统", Status: 1,
	//})
	//orm.Create(&model.Shipping{
	//	Title: "店家配送", Key: "shop", Intro: "基于支付宝提供的支付系统", Status: 1,
	//})
	//orm.Create(&model.Shipping{
	//	Title: "顺丰", Key: "shunfeng", Intro: "基于银联提供的支付系统", Status: 1,
	//})
	//orm.Create(&model.Shipping{
	//	Title: "EMS", Key: "ems", Intro: "基于银联提供的支付系统", Status: 1,
	//})
	//log.Println("table <shipping> is seeded!")
	//
	//orm.Create(&model.ShippingStatus{
	//	Title: "配送错误",
	//})
	//orm.Create(&model.ShippingStatus{
	//	Title: "未配送",
	//})
	//orm.Create(&model.ShippingStatus{
	//	Title: "已配送",
	//})
	//orm.Create(&model.ShippingStatus{
	//	Title: "已收货",
	//})
	//
	//log.Println("table <shipping_status> is seeded!")
	//
	//
	//orm.Create(&model.OrderStatus{
	//	Title: "订单错误",
	//})
	//orm.Create(&model.OrderStatus{
	//	Title: "确认",
	//})
	//orm.Create(&model.OrderStatus{
	//	Title: "完成",
	//})
	//orm.Create(&model.OrderStatus{
	//	Title: "取消",
	//})
	//orm.Create(&model.OrderStatus{
	//	Title: "删除",
	//})
	//log.Println("table <order_status> is seeded!")

	// 3.1 Category
	//orm.Create(&model.Category{
	//	Name: "未分类", ParentId: 0,
	//})
	//orm.Create(&model.Category{
	//	Name: "图书", ParentId: 0,
	//})
	//orm.Create(&model.Category{
	//	Name: "电脑", ParentId: 0,
	//})
	//log.Println("table <Category> is seeded!")

	// 3.2 Product
	//orm.Create(&model.Product{
	//	Name: "图书 A", CategoryID: 2, Upc: "177173627181192",
	//})
	//orm.Create(&model.Product{
	//	Name: "图书 B", CategoryID: 2, Upc: "177173627181193",
	//})
	//orm.Create(&model.Product{
	//	Name: "图书 C", CategoryID: 2, Upc: "177173627181194",
	//})
	//orm.Create(&model.Product{
	//	Name: "图书 D", CategoryID: 2, Upc: "177173627181195",
	//})
	//orm.Create(&model.Product{
	//	Name: "图书 E", CategoryID: 2, Upc: "177173627181196",
	//})
	//orm.Create(&model.Product{
	//	Name: "图书 F", CategoryID: 2, Upc: "177173627181197",
	//})
	//orm.Create(&model.Product{
	//	Name: "电脑 A", CategoryID: 3, Upc: "177173627181166",
	//})
	//orm.Create(&model.Product{
	//	Name: "电脑 B", CategoryID: 3, Upc: "177173627181116",
	//})
	//orm.Create(&model.Product{
	//	Name: "电脑 C", CategoryID: 3, Upc: "177173627181126",
	//})
	//orm.Create(&model.Product{
	//	Name: "电脑 D", CategoryID: 3, Upc: "177173627181136",
	//})
	//orm.Create(&model.Product{
	//	Name: "电脑 E", CategoryID: 3, Upc: "177173627181146",
	//})
	//orm.Create(&model.Product{
	//	Name: "电脑 F", CategoryID: 3, Upc: "177173627181156",
	//})
	//log.Println("table <Product> is seeded!")
	//
	//orm.Create(&model.User{
	//	User: "root",
	//})
	//log.Println("table <Product> is seeded!")
}

