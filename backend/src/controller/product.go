package controller

import (
	"config"
	"github.com/gin-gonic/gin"
	"model"
	"net/http"
	"strconv"
	"time"
)


//复制
func ProductCopy(c *gin.Context) {
	// 先获取源商品信息（全部）
	IDstr := c.DefaultQuery("ID", "")
	if IDstr == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "请输入资源 ID",
		})
		return
	}
	ID, _ := strconv.Atoi(IDstr)

	// 2 得到模型
	src := model.Product{} // 源产品模型
	src.ID = uint(ID)
	orm.Find(&src)
	// 包含关联数据。本例中主要是产品属性ProductAttr
	//orm.Model(&src).Related(&src.ProductAttrs)

	// 拷贝新产品
	dst := src// 目标新产品
	// 清理相关的标识属性
	dst.ID = 0
	dst.Upc = strconv.Itoa(int(time.Now().UnixNano()))
	// 生成新产品
	orm.Create(&dst)

	// 拷贝关联数据。
	orm.Model(&src).Related(&src.ProductAttrs)
	// 遍历
	for _, pa := range src.ProductAttrs {
		dstPa := model.ProductAttr{}
		dstPa.AttrID = pa.AttrID
		dstPa.Value = pa.Value
		dstPa.ProductID = dst.ID
		// 产品某个属性的复制
		orm.Create(&dstPa)
	}

	orm.Model(&dst).
		Related(&dst.ProductAttrs).
		Related(&dst.Category)

	// 保证属性与值的映射关系
	dst.AttrValue = map[uint]string{}
	for _, pa := range dst.ProductAttrs {
		dst.AttrValue[pa.AttrID] = pa.Value
	}
	dst.ProductAttrs = nil

	// 从产品的属性类型考虑，考虑全部的属性
	if !orm.Model(&dst).Related(&dst.AttrType).RecordNotFound() {
		// 存在类型
		//根据类型确认全部的属性
		orm.Model(&dst.AttrType).Related(&dst.AttrType.AttrGroups)
		// 通过group找到对应的全部属性
		for ii, ag := range dst.AttrType.AttrGroups {
			// 查询关联的属性
			orm.Model(&ag).Related(&dst.AttrType.AttrGroups[ii].Attrs)
			for _, a := range dst.AttrType.AttrGroups[ii].Attrs {
				// 得到该产品的全部属性
				if _, exists := dst.AttrValue[a.ID]; !exists {
					dst.AttrValue[a.ID] = ""
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": dst,
	})
}

//更新
func ProductUpdate(c *gin.Context) {
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
	m := model.Product{}
	m.ID = uint(ID)
	//orm.Find(&m)

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

	// 更新produc-attr 表数据。先获取，该产品应该具有哪些属性。形成对应的记录。若存在更新的值，将其更新。
	if !orm.Model(&m).Related(&m.AttrType).RecordNotFound() {
		// 存在类型
		//根据类型确认全部的属性
		orm.Model(&m.AttrType).Related(&m.AttrType.AttrGroups)
		// 通过group找到对应的全部属性
		for i, ag := range m.AttrType.AttrGroups {
			// 查询关联的属性
			orm.Model(&ag).Related(&m.AttrType.AttrGroups[i].Attrs)
			for _, a := range m.AttrType.AttrGroups[i].Attrs {
				// 根据属性选择更新或插入
				pa := model.ProductAttr{}
				if orm.Model(&model.ProductAttr{}).
					Where("product_id=? AND attr_id=?", m.ID, a.ID).
					Find(&pa).RecordNotFound() {
					// 不存在
					pa.Value = m.AttrValue[a.ID]
					pa.ProductID = m.ID
					pa.AttrID = a.ID
					orm.Create(&pa)
				} else {
					// 已经存在
					pa.Value = m.AttrValue[a.ID]// 新设置的值
					orm.Save(&pa)
				}
			}
		}
	}

	// 更新image表
	for i, img := range m.UploadedImage {
		// a/b/ab000000.jpg
		image := model.Image{}
		image.ProductID = m.ID
		image.Host = config.App["IMAGE_HOST"]
		image.Image = string(img[0]) + "/" + string(img[1]) + "/" + img
		image.ImageSmall = string(m.UploadedImageSmall[i][0]) + "/" + string(m.UploadedImageSmall[i][1]) + "/" + m.UploadedImageSmall[i]
		image.ImageBig = string(m.UploadedImageBig[i][0]) + "/" + string(m.UploadedImageBig[i][1]) + "/" + m.UploadedImageBig[i]
		orm.Create(&image)
	}

	// 5 正确响应
	// 查询相应的关联数据
	category := model.Category{}
	category.ID = m.CategoryID
	orm.Find(&category)
	m.Category = category
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": m,
	})
}

//添加
func ProductCreate(c *gin.Context) {
	// 1 得到模型
	m := model.Product{}
	// 2 从请求数据中绑定对象
	bindErr := c.ShouldBind(&m)
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": bindErr.Error(),
		})
		return
	}
	// 特定数据的设置
	if "" == m.Upc { // 用户未填写
		m.Upc = strconv.Itoa(int(time.Now().UnixNano()))
	}
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
	// product 处理成功，处理关联数据
	// 更新produc-attr 表数据。先获取，该产品应该具有哪些属性。形成对应的记录。若存在更新的值，将其更新。
	if !orm.Model(&m).Related(&m.AttrType).RecordNotFound() {
		// 存在类型
		//根据类型确认全部的属性
		orm.Model(&m.AttrType).Related(&m.AttrType.AttrGroups)
		// 通过group找到对应的全部属性
		for i, ag := range m.AttrType.AttrGroups {
			// 查询关联的属性
			orm.Model(&ag).Related(&m.AttrType.AttrGroups[i].Attrs)
			//log.Println(m.AttrType.AttrGroups[i].Attrs)
			for _, a := range m.AttrType.AttrGroups[i].Attrs {
				// 根据属性选择更新或插入
				pa := model.ProductAttr{}
				if orm.Model(&model.ProductAttr{}).
					Where("product_id=? AND attr_id=?", m.ID, a.ID).
					Find(&pa).RecordNotFound() {
					// 不存在
					pa.Value = m.AttrValue[a.ID]
					pa.ProductID = m.ID
					pa.AttrID = a.ID
					orm.Create(&pa)
				} else {
					// 已经存在
					pa.Value = m.AttrValue[a.ID]// 新设置的值
					orm.Save(&pa)
				}
			}
		}
	}


	//查询相应的关联数据
	category := model.Category{}
	category.ID = m.CategoryID
	orm.Find(&category)
	m.Category = category
	// 4 正确响应
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": m,
	})

}

// 删除
func ProductDelete(c *gin.Context) {
	ID := c.DefaultQuery("ID", "")
	if "" == ID {
		c.JSON(http.StatusOK, gin.H{
			"error": "未指定资源 ID",
		})
		return
	}

	// 确定模型
	m := &model.Product{}
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

// 产品列表接口
func ProductList(c *gin.Context) {

	// # 搜索
	condStr := ""
	condParams := []string{}
	// 确定搜索条件
	filterName := c.DefaultQuery("filterName", "") // Name
	if filterName != "" {
		// 需要搜索
		condStr = "name like ?"
		condParams = append(condParams, filterName + "%")
	}
	// 其余的条件

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
	orm.Model(&model.Product{}).
		Where(condStr, condParams).
		Count(&total)

	// 计算 偏移
	offset := (currentPage - 1) * pageSize

	// # 获取数据
	products := []model.Product{}
	orm.
		Where(condStr, condParams).
		Order(orderStr).
		Limit(pageSize).
		Offset(offset).
		Find(&products)
	// 遍历每个Product 得到对应的关联Category
	for i, _ := range products {
		orm.Model(&products[i]).Related(&products[i].Category)
		// 查询属性
		products[i].AttrValue = map[uint]string{}
		if !orm.Model(&products[i]).Related(&products[i].ProductAttrs).RecordNotFound() {
			for _, pa := range products[i].ProductAttrs {
				products[i].AttrValue[pa.AttrID] = pa.Value
			}
			products[i].ProductAttrs = nil
		}

		// 从产品的属性类型考虑，考虑全部的属性
		if !orm.Model(&products[i]).Related(&products[i].AttrType).RecordNotFound() {
			// 存在类型
			//根据类型确认全部的属性
			orm.Model(&products[i].AttrType).Related(&products[i].AttrType.AttrGroups)
			// 通过group找到对应的全部属性
			for ii, ag := range products[i].AttrType.AttrGroups {
				// 查询关联的属性
				orm.Model(&ag).Related(&products[i].AttrType.AttrGroups[ii].Attrs)
				for _, a := range products[i].AttrType.AttrGroups[ii].Attrs {
					// 得到该产品的全部属性
					if _, exists := products[i].AttrValue[a.ID]; !exists {
						products[i].AttrValue[a.ID] = ""
					}
				}
			}
		}

		// 查询关联图像
		products[i].Images = []model.Image{}
		orm.Model(&products[i]).Related(&products[i].Images)
	}

	// 响应
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": products,
		"pager": map[string]int {
			"currentPage": currentPage,
			"pageSize": pageSize,
			"total": total,
		},
	})

}

