package middleware

import (
	"config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"model"
	"net/http"
)

//获取用户权限中间件

func Pri(c *gin.Context) {

	// 利用userName的用户的权限列表
	//初始化 gorm
	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return config.App["DB_TABLE_PREFIX"] + defaultTableName
	}
	// 连接数据库
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
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userName, _ := c.Get("userName")
	user := model.User{}
	if orm.Where("user = ?", userName).Find(&user).RecordNotFound() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}
	pris := []string{}

	// 利用user得到权限列表
	role := model.Role{}
	orm.Model(&user).Related(&role, "Role")
	if role.ID != 0 { // 存在对应的角色
		// 利用角色获取权限
		cps := []model.Privilege{}
		orm.Model(&role).Related(&cps, "Privileges")
		// 拼凑一个权限可以的切片

		for _, p := range cps {
			pris = append(pris, p.Key)
		}
	}
	// 记录当前用户的权限
	c.Set("pris", pris)
}