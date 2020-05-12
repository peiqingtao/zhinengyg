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
	// 3 迁移（利用模型迁移）(migrate
	orm.AutoMigrate(&model.Product{})
	log.Println("table <product> is migrated!")
	// 插入Product测试数据,(Seed)

	//// 3.2 Category 的 migrate
	//orm.AutoMigrate(&model.Category{})
	//log.Println("table <Category> is created!")

	orm.AutoMigrate(&model.User{})
	log.Println("table <user> is created!")

	orm.AutoMigrate(&model.Role{})
	log.Println("table <role> is created!")

	orm.AutoMigrate(&model.Privilege{})
	log.Println("table <privilege> is created!")

	orm.AutoMigrate(&model.RolePrivilege{})
	log.Println("table <role_privilege> is created!")

	orm.AutoMigrate(&model.AttrType{})
	log.Println("table <attr_type> is created!")

	orm.AutoMigrate(&model.AttrGroup{})
	log.Println("table <attr_group> is created!")

	orm.AutoMigrate(&model.Attr{})
	log.Println("table <attr> is created!")

	orm.AutoMigrate(&model.ProductAttr{})
	log.Println("table <product_attr> is created!")

	orm.AutoMigrate(&model.Group{})
	log.Println("table <group> is created!")

	orm.AutoMigrate(&model.GroupAttr{})
	log.Println("table <group_attr> is created!")

	orm.AutoMigrate(&model.Image{})
	log.Println("table <image> is created!")

	orm.AutoMigrate(&model.Cart{})
	log.Println("table <cart> is created!")

	orm.AutoMigrate(&model.Address{})
	log.Println("table <address> is created!")

	orm.AutoMigrate(&model.Order{})
	log.Println("table <order> is created!")

	orm.AutoMigrate(&model.Payment{})
	log.Println("table <payment> is created!")
	orm.AutoMigrate(&model.PaymentStatus{})
	log.Println("table <payment_status> is created!")

	orm.AutoMigrate(&model.Shipping{})
	log.Println("table <shipping> is created!")
	orm.AutoMigrate(&model.ShippingStatus{})
	log.Println("table <shipping_status> is created!")

	orm.AutoMigrate(&model.OrderStatus{})
	log.Println("table <order_status> is created!")

	orm.AutoMigrate(&model.OrderProduct{})
	log.Println("table <order_product> is created!")

	orm.AutoMigrate(&model.Quantity{})
	log.Println("table <quantity> is created!")


}
