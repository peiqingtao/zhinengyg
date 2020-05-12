package controller

// 产品属性 控制器动作列表

import (
	"github.com/gin-gonic/gin"
	"model"
	"net/http"
	"strconv"
)

//更新
func ProductAttrUpdate(c *gin.Context) {
	// 1 确定资源 ID
	IDstr := c.DefaultQuery("ID", "")
	if IDstr == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "请输入资源 ID",
		})
		return
	}
	ID, _ := strconv.Atoi(IDstr)

	// 2 得到模型
	m := model.ProductAttr{}
	m.ID = uint(ID)
	orm.Find(&m)

	// 3 从请求数据中绑定对象
	bindErr := c.ShouldBind(&m)
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": bindErr.Error(),
		})
		return
	}
	// 特定数据的设置

	// 4 更新
	//将关联的存储临时取消
	orm.
		Set("gorm:save_associations", false).
		Save(&m)
	// 插入有错误
	if orm.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": orm.Error.Error(),
		})
		return
	}

	// 5 正确响应
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": m,
	})
}

//添加
func ProductAttrCreate(c *gin.Context) {
	// 1 得到模型
	m := model.ProductAttr{}
	// 2 从请求数据中绑定对象
	bindErr := c.ShouldBind(&m)
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": bindErr.Error(),
		})
		return
	}
	// 特定数据的设置

	// 3 插入
	//将关联的存储临时取消
	orm.
	    Set("gorm:save_associations", false).
		Create(&m)
	// 插入有错误
	if orm.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": orm.Error.Error(),
		})
		return
	}
	//查询相应的关联数据

	// 4 正确响应
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": m,
	})

}

// 删除
func ProductAttrDelete(c *gin.Context) {
	ID := c.DefaultQuery("ID", "")
	if "" == ID {
		c.JSON(http.StatusOK, gin.H{
			"error": "未指定资源 ID",
		})
		return
	}

	// 确定模型
	m := &model.ProductAttr{}
	id, _ := strconv.Atoi(ID)
	m.ID = uint(id)

	// 删除
	orm.Delete(&m)
	if orm.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": orm.Error.Error(),
		})
		return
	}

	// 无错误的响应
	c.JSON(http.StatusOK, gin.H{
		"error": "",
	})
}

// 列表接口
func ProductAttrList(c *gin.Context) {

	// # 搜索
	condStr := ""
	condParams := []string{}
	// 确定搜索条件

	// # 排序
	// 排序参数
	orderStr := ""
	sortProp := c.DefaultQuery("sortProp", "")
	sortOrder := c.DefaultQuery("sortOrder", "")// 方式 ，ascending, descending
	if sortProp != "" && sortOrder != "" { // 请求了某个字段的排序
		sortMethod := "ASC"
		if "descending" == sortOrder {
			sortMethod = "DESC"
		}
		orderStr = sortProp + " " + sortMethod // name asc
	}


	// # 翻页
	currentPageStr := c.DefaultQuery("currentPage", "1") // 当前页，默认为 1
	pageSizeStr := c.DefaultQuery("pageSize", "5")
	// 当前页
	currentPage, pageErr := strconv.Atoi(currentPageStr)
	if pageErr != nil {
		currentPage = 1
	}
	// 每页size
	pageSize, pageErr := strconv.Atoi(pageSizeStr)
	if pageErr != nil {
		pageSize = 5
	}

	// 获取总记录数
	total := 0
	orm.Model(&model.ProductAttr{}).
		Where(condStr, condParams).
		Count(&total)

	// 计算 偏移
	offset := (currentPage - 1) * pageSize

	// # 获取数据
	ms := []model.ProductAttr{}
	orm.
		Where(condStr, condParams).
		Order(orderStr).
		Limit(pageSize).
		Offset(offset).
		Find(&ms)

	// 响应
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": ms,
		"pager": map[string]int {
			"currentPage": currentPage,
			"pageSize": pageSize,
			"total": total,
		},
	})

}

