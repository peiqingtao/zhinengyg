package router

import (
	"controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RouterInit() *gin.Engine {
	// 1.初始化路由引擎对象
	r := gin.Default()
	r.Use(cors.Default())
	// 2.定义路由，以及对应的动作处理函数
	r.GET("/ping", controller.Ping)
	//很多的路由对应动作...
	r.GET("/category-tree", controller.CategoryTree)
	r.POST("/category", controller.CategoryAdd)
	r.DELETE("/category", controller.CategoryDelete)

	// 返回路由引擎对象
	return r
}
