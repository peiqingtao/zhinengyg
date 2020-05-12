package controller

// 用户 控制器动作列表

import (
	"config"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"model"
	"net/http"
	"strconv"
	"strings"
	"time"
)


// 认证
func UserAuth(c *gin.Context) {
	// 先去校验验证码
	postCode := c.PostForm("Code")
	if postCode == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "请填写验证码",
		})
		return
	}
	i := strings.Index(c.Request.RemoteAddr, "]")
	key := fmt.Sprintf("%x", md5.Sum([]byte(c.Request.RemoteAddr[1:i] + c.Request.Header["User-Agent"][0])))
	result, err := Rds.Do("get", "code_" + key)
	code := string(result.([]byte))
	if postCode != code {
		c.JSON(http.StatusOK, gin.H{
			"error": "验证码错误",
		})
		return
	}
	// 验证通过，删除验证码
	Rds.Do("DEL", "code_" + key)



	user := model.User{}
	// 先通过用户名获取用户信息，再比较密码是否正确
	postUser := c.PostForm("User")
	if orm.Where("user=?", postUser).Find(&user).RecordNotFound() {
		// 没有该用户
		c.JSON(http.StatusOK, gin.H{
			"error": "用户或密码错误，用户名",
		})
		return
	}

	// 检测密码
	postPassword := c.PostForm("Password")
	pwdFunc := hmac.New(sha256.New, []byte(user.PasswordSalt)) // 生成一个 hs函数，使用数据表中存储的salt来处理
	pwdFunc.Write([]byte(postPassword)) // 写入待生成摘要的内容，原密码
	if user.Password != fmt.Sprintf("%x", pwdFunc.Sum(nil)) {// 计算摘要值。Sum，摘要（CheckSum）
		// 密码错误
		c.JSON(http.StatusOK, gin.H{
			"error": "用户或密码错误，密码",
		})
		return
	}

	// 认证通过
	// 生成 JWT 格式的 Token
	// 生成签名的key
	mySigningKey := []byte(config.App["SECRET"])

	// token中包含如下的数据
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 7200, // 有效期 当前时间戳+有效期
		Issuer:    "Backend", // backend 发行
		Audience: user.User, // 用户名字
	}
	//生成 token (先得到token构建器，在生成 token）
	tokenBuilder := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenBuilder.SignedString(mySigningKey)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "Token 生成失败",
		})
		return
	}

	// 正确响应
	user.PasswordSalt, user.Password = "", ""
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"user": user,
		"token": token,
	})

}

//更新
func UserUpdate(c *gin.Context) {
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
	m := model.User{}
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
	log.Println(m)

	// 4 更新
	//将关联的存储临时取消
	orm.Model(&m).
		Set("gorm:save_associations", false).
		Update(m)
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
func UserCreate(c *gin.Context) {
	// 1 得到模型
	m := model.User{}
	// 2 从请求数据中绑定对象
	bindErr := c.ShouldBind(&m)
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": bindErr.Error(),
		})
		return
	}
	// 特定数据的设置
	// 生成salt
	saltChars := "abcdefghijklnmopqrstuvwxyz1234567890!@#$%^&*()" //salt 可能的字符串
	saltLen := 6 // slat 长度
	salt := ""
	for i:=0; i<saltLen; i++ {
		index := rand.Int31n(int32(len(saltChars)))
		salt += string(saltChars[index])
	}
	m.PasswordSalt = salt
	// 为密码做信息摘要（加密），hmac sha256
	pwdFunc := hmac.New(sha256.New, []byte(salt)) // 生成一个 hs函数
	pwdFunc.Write([]byte(m.Password)) // 写入待生成摘要的内容，原密码
	m.Password = fmt.Sprintf("%x", pwdFunc.Sum(nil)) // 计算摘要值。Sum，摘要（CheckSum）

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
func UserDelete(c *gin.Context) {
	ID := c.DefaultQuery("ID", "")
	if "" == ID {
		c.JSON(http.StatusOK, gin.H{
			"error": "未指定资源 ID",
		})
		return
	}

	// 确定模型
	m := &model.User{}
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
func UserList(c *gin.Context) {

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
	orm.Model(&model.User{}).
		Where(condStr, condParams).
		Count(&total)

	// 计算 偏移
	offset := (currentPage - 1) * pageSize

	// # 获取数据
	ms := []model.User{}
	orm.
		Where(condStr, condParams).
		Order(orderStr).
		Limit(pageSize).
		Offset(offset).
		Find(&ms)
	for i, _ := range ms {
		ms[i].PasswordSalt = ""
		ms[i].Password = ""
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

