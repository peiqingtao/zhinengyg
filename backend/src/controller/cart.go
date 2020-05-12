package controller

import (
	"github.com/gin-gonic/gin"
	"model"
	"net/http"
)

func CartProduct(c *gin.Context) {
	// # 搜索
	// 确定搜索条件
	filterIDs := c.QueryArray("filterIDs[]")
	if len(filterIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"error": "产品ID不存在",
		})
		return
	}

	// 拼凑 condStr
	condStr := "id in (?)"
	// # 获取数据
	products := []model.Product{}
	orm.
		Where(condStr, filterIDs).
		Find(&products)
	// 遍历每个Product 得到对应的关联Category

	for i, _ := range products {
		// 查询关联图像
		products[i].Images = []model.Image{}
		orm.Model(&products[i]).Related(&products[i].Images)

		// 产品的型号信息
		productModel(&products[i])
	}

	// 响应
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": products,
	})
}
