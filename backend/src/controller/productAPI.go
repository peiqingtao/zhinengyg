package controller

import (
	"github.com/gin-gonic/gin"
	"model"
	"net/http"
	"strconv"
	"strings"
)



//获取产品型号信息
func productModel(product *model.Product) {
	// 如果组存在，
	if product.GroupID == 0 {
		return
	}

	//先获取组
	orm.Model(&product).Related(&product.Group)

	// 再获取组的差异属性列表
	orm.Model(&product.Group).Related(&product.Group.Attrs, "Attrs")

	if len(product.Group.Attrs) == 0 {
		return
	}

	// 检测是否存在差异属性
	aids := []uint{}
		// 获取当前商品的该属性值
	for _, a := range product.Group.Attrs {
		aids = append(aids, a.ID)
	}

	pas := []model.ProductAttr{}
	orm.Where("product_id=? AND attr_id in (?)", product.ID, aids).Find(&pas)

	values := []string{}
	for _, pa := range pas {
		values = append(values, pa.Value)
	}

	product.ModelInfo = strings.Join(values, "|")
}

//获取产品信息
func ProductInfo(c *gin.Context) {
	IDstr := c.DefaultQuery("ID", "")
	if IDstr == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "请输入资源 ID",
		})
		return
	}
	ID, _ := strconv.Atoi(IDstr)

	// 2 得到产品模型
	product := model.Product{}
	product.ID = uint(ID)
	orm.Find(&product)

	// 3 获取关联数据
	// 分类，图像
	orm.Model(&product).Related(&product.Category)

	// 查询关联图像
	product.Images = []model.Image{}
	orm.Model(&product).Related(&product.Images)

	// 型号信息
	productModel(&product)

	// 关联的组
	orm.Model(&product).Related(&product.Group)
	product.Group.Products = []model.Product{}
	if product.Group.ID != 0 { // 存在分组
		//查找组内产品
		orm.Model(&product.Group).Related(&product.Group.Products)
		for i, _ := range product.Group.Products {
			productModel(&product.Group.Products[i])
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": product,
	})

}

//获取推荐商品
func ProductPromote(c *gin.Context) {
	// # 搜索
	condStr := ""
	condParams := []string{}
	// 确定搜索条件
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
