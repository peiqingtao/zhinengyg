package controller

// 属性分组 控制器动作列表

import (
	"github.com/gin-gonic/gin"
	"model"
	"net/http"
	"strconv"
	"strings"
)

//更新
func AttrGroupUpdate(c *gin.Context) {
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
	m := model.AttrGroup{}
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
func AttrGroupCreate(c *gin.Context) {
	// 1 得到模型
	m := model.AttrGroup{}
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
func AttrGroupDelete(c *gin.Context) {
	ID := c.DefaultQuery("ID", "")
	if "" == ID {
		c.JSON(http.StatusOK, gin.H{
			"error": "未指定资源 ID",
		})
		return
	}

	// 确定模型
	m := &model.AttrGroup{}
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
func AttrGroupList(c *gin.Context) {

	// # 搜索
	condStr := ""
	cond := []string{}
	condParams := []string{}
	// 确定搜索条件
	filterAttrTypeID := c.DefaultQuery("filterAttrTypeID", "") // Name
	if filterAttrTypeID != "" {
		cond = append(cond, "attr_type_id = ?")
		condParams = append(condParams, filterAttrTypeID)
	}

	// 拼凑 condStr
	condStr = strings.Join(cond, " AND ")

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
	orm.Model(&model.AttrGroup{}).
		Where(condStr, condParams).
		Count(&total)

	// 计算 偏移
	offset := (currentPage - 1) * pageSize

	// # 获取数据
	ms := []model.AttrGroup{}
	orm.
		Where(condStr, condParams).
		Order(orderStr).
		Limit(pageSize).
		Offset(offset).
		Find(&ms)


	// 关联查询
	withAttr := c.DefaultQuery("withAttr", "")
	for i, m := range ms {
		// 关联类型
		orm.Model(&m).Related(&ms[i].AttrType)

		// 查询关联的属性
		if withAttr == "yes" {
			orm.Model(&m).Related(&ms[i].Attrs)
		}
	}


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

