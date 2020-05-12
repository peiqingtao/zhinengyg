package controller

// 分组 控制器动作列表

import (
	"github.com/gin-gonic/gin"
	"model"
	"net/http"
	"strconv"
)

//更新
func GroupUpdate(c *gin.Context) {
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
	m := model.Group{}
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
	if len(m.CheckedProductID) > 0 {
		// 用户选择了所属产品
		orm.Model(&model.Product{}).Where("id in (?)", m.CheckedProductID).Update("group_id", m.ID)
		orm.Model(&model.Product{}).Where("id not in (?) AND group_id = ?", m.CheckedProductID, m.ID).Update("group_id", 0)
	}
	// 处理差异属性
	if len(m.CheckedAttrID) > 0 {
		// 构建一个 []Attr
		as := []model.Attr{}
		for _, aID := range m.CheckedAttrID {
			a := model.Attr{}
			a.ID = aID
			as = append(as, a)
		}
		// 利用多堆多替换
		orm.Model(&m).Association("Attrs").Replace(as)
	}


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
func GroupCreate(c *gin.Context) {
	// 1 得到模型
	m := model.Group{}
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
	// 特定数据的设置
	if len(m.CheckedProductID) > 0 {
		// 用户选择了所属产品
		orm.Model(&model.Product{}).Where("id in (?)", m.CheckedProductID).Update("group_id", m.ID)
	}

	// 4 正确响应
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": m,
	})

}

// 删除
func GroupDelete(c *gin.Context) {
	ID := c.DefaultQuery("ID", "")
	if "" == ID {
		c.JSON(http.StatusOK, gin.H{
			"error": "未指定资源 ID",
		})
		return
	}

	// 确定模型
	m := &model.Group{}
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
func GroupList(c *gin.Context) {

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
	orm.Model(&model.Group{}).
		Where(condStr, condParams).
		Count(&total)

	// 计算 偏移
	offset := (currentPage - 1) * pageSize

	// # 获取数据
	ms := []model.Group{}
	orm.
		Where(condStr, condParams).
		Order(orderStr).
		Limit(pageSize).
		Offset(offset).
		Find(&ms)

	// 获取关联数据
	for i, m := range ms {
		// 关联产品
		orm.Model(&m).Related(&ms[i].Products)
		// 遍历全部的产品，将ID做集合
		ms[i].CheckedProductID = []uint{}
		for _, p := range ms[i].Products {
			ms[i].CheckedProductID = append(ms[i].CheckedProductID, p.ID)
		}

		//关联差异属性
		//orm.Model(&m).Related(&ms[i].Attrs, "Attrs")
		orm.Model(&m).Association("Attrs").Find(&ms[i].Attrs)
		// 遍历全部的产品，将ID做集合
		ms[i].CheckedAttrID = []uint{}
		for _, a := range ms[i].Attrs {
			ms[i].CheckedAttrID = append(ms[i].CheckedAttrID, a.ID)
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

