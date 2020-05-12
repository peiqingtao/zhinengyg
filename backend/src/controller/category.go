package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"model"
	"strconv"
)

// 负责 分类相关操作的函数集合文件

// 添加分类
func CategoryAdd(c *gin.Context) {
	// 连接数据库
	db, dberr := gorm.Open("mysql", "projectAUser:hellokang@/projectA?charset=utf8mb4&loc=Local")
	if dberr != nil {
		log.Println(dberr)
		return
	}
	defer db.Close()

	// 得到一个Category模型对象
	category := model.Category{}
	// 从请求中解析数据，保证请求数据的格式与模型一致
	err := c.ShouldBind(&category)
	if  err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 自动解析绑定成功
	// 完成数据校验
	// 数据入库
	db.Create(&category)

	// 完成响应
	c.JSON(200, gin.H{
		"error": "",
		"data": category,
	})
}

// 分类树
func CategoryTree(c *gin.Context) {
	// 连接数据库
	db, err := gorm.Open("mysql", "projectAUser:hellokang@/projectA?charset=utf8mb4")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	// 查询全部的分类
	categories := []model.Category{}
	db.Find(&categories)

	// 响应 JSON
	c.JSON(200, gin.H{
		"error": "",
		"data": categories,
		})
}

// 分类树
//func CategoryTree(c *gin.Context) {
//	//连接数据库，获取全部的分类内容
//	config := map[string]string{
//		"username": "projectAUser",
//		"password": "hellokang",
//		"host": "127.0.0.1",
//		"port": "3306",
//		"dbname": "projectA",
//		"collation": "utf8mb4_general_ci",
//	}
//	db, err := dao.NewDao(config)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	// 查询分类的全部数据
//	rows, err := db.Table("a_categories").Rows()
//	log.Println(rows, err)
//	if err == nil {
//		// 响应结果
//		c.JSON(200, gin.H{
//			"error": "",
//			"data": rows,
//		})
//	} else {
//		c.JSON(200, gin.H{
//			"error": err.Error(),
//		})
//	}
//
//}

//动作 2
func CategoryDelete(c * gin.Context) {
	// 连接数据库
	db, dberr := gorm.Open("mysql", "projectAUser:hellokang@/projectA?charset=utf8mb4&loc=Local")
	if dberr != nil {
		log.Println(dberr)
		return
	}
	defer db.Close()

	// 需要获取删除的ID
	ID := c.Query("ID")
	// 构建模型对象
	category := model.Category{}
	id, _ := strconv.Atoi(ID)
	category.ID = uint(id)

	// 删除
	db.Delete(&category)
	c.JSON(200, gin.H{
		"error": "",
	})
}