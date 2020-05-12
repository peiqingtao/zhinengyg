package controller

import (
	"bytes"
	"config"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"model"
	"net/http"
	"time"
)



func member(c *gin.Context) *model.User {
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		return nil
	}

	token := string(bytes.Replace([]byte(authorization), []byte("Bearer "), []byte(""), -1))
	if token == "" {
		return nil
	}

	// 校验 token 是否被篡改
	tokenObj, parseErr := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.App["SECRET"]), nil
	})
	if parseErr != nil { // token 语法错误
		return nil
	}
	// 判断解析结果，验证通过 使用 .Valid 表示
	if !tokenObj.Valid { // 签名错误
		return nil
	}

	// token 验证通过
	userName:=""
	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok {
		userName = claims["aud"].(string)
	}
	// 校验通过，获取用户名称，利用用户名称获取用户全部信息
	user := model.User{}
	orm.Where("user=?", userName).Find(&user)
	return &user
}

// 地址添加
func MemberAddressAdd(c *gin.Context) {
	user := member(c)
	if user == nil {
		// 未认证通过
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	address := model.Address{}
	// 2 从请求数据中绑定对象
	bindErr := c.ShouldBind(&address)
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": bindErr.Error(),
		})
		return
	}

	address.UserID = user.ID
	// 如果当前的地址为默认，应该将其他的默认地址给去掉。
	log.Println(address.IsDefault)
	if address.IsDefault {
		orm.Model(&model.Address{}).Where("user_id=? AND is_default=1", user.ID).Update("is_default", 0)
	}

	orm.Create(&address)
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": address,
	})

}

// 地址列表
func MemberAddressList (c *gin.Context) {
	user := member(c)
	if user == nil {
		// 未认证通过
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// 查询
	as := []model.Address{}
	orm.Where("user_id=?", user.ID).
		Order("is_default desc").
		Find(&as)
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": as,
	})
}

func MemberCartSet(c *gin.Context) {

	user := member(c)
	if user == nil {
		// 未认证通过
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Auth failed!",
		})
		return
	}

	// 执行同步。
	// 将 浏览器端端购物车的数据与会员已有的购物车的数据合并起来。
	// 浏览器携带的
	tempCart := c.PostForm("cart")
	tempCartProducts := []struct{
		ProductID uint  `json:"productID"`
		BuyQuantity int `json:"buyQuantity"`
	}{}
	json.Unmarshal([]byte(tempCart), &tempCartProducts)


	// 同步完成， 存储到cart表中。
	//json 编码
	content, _ := json.Marshal(tempCartProducts)

	// 更新
	memberCart := model.Cart{}
	orm.Where("user_id=?", user.ID).Find(&memberCart)
	memberCart.Content = string(content)
	orm.Save(&memberCart)

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": string(content),
	})
}


func MemberCartSync(c *gin.Context) {

	user := member(c)
	if user == nil {
		// 未认证通过
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Auth failed!",
		})
		return
	}

	// 执行同步。
	// 将 浏览器端端购物车的数据与会员已有的购物车的数据合并起来。
	// 浏览器携带的
	tempCart := c.PostForm("cart")
	tempCartProducts := []struct{
		ProductID uint  `json:"productID"`
		BuyQuantity int `json:"buyQuantity"`
	}{}
	json.Unmarshal([]byte(tempCart), &tempCartProducts)
	// 服务器端已有的
	memberCart := model.Cart{}
	orm.Where("user_id=?", user.ID).Find(&memberCart)
	memberCartProducts := []struct{
		ProductID uint  `json:"productID"`
		BuyQuantity int `json:"buyQuantity"`
	}{}
	if memberCart.Content != "" {
		json.Unmarshal([]byte(memberCart.Content), &memberCartProducts)
	}

	//同步，遍历全部的浏览器端携带的
	for _, tp := range tempCartProducts {
		// 判断是否存在与服务器端的购物车产品中
		exists := false // 假设不存在
		for _, mp := range memberCartProducts {
			if tp.ProductID == mp.ProductID {
				exists = true// 存在
				break
			}
		}
		// 判断是否存在，若不存在，则增加商品即可
		if !exists {
			memberCartProducts = append(memberCartProducts, tp)
		}
	}

	// 同步完成， 存储到cart表中。
	//json 编码
	content, _ := json.Marshal(memberCartProducts)

	if  memberCart.ID == 0 {
		// 添加购物车数据
		memberCart.UserID = user.ID
		memberCart.Content = string(content)
		orm.Create(&memberCart)
	} else  {
		// 更新
		memberCart.Content = string(content)
		orm.Save(&memberCart)
	}

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": string(content),
	})
}

func MemberCart(c *gin.Context) {
	user := member(c)
	if user == nil {
		// 未认证通过
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Auth failed!",
		})
		return
	}

	// 查找会员的购物车信息
	cart := model.Cart{}
	orm.Where("user_id=?", user.ID).Find(&cart)
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": cart.Content,
	})
}

func MemberAuth(c *gin.Context) {

	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		// 没有 authorization 头
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "header Authorization not found",
		})
		return
	}
	token := string(bytes.Replace([]byte(authorization), []byte("Bearer "), []byte(""), -1))

	if token == "" {
		// 没有 Token
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "token not found",
		})
		c.Abort()
		return
	}

	// 校验 token 是否被篡改
	tokenObj, parseErr := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.App["SECRET"]), nil
	})
	if parseErr != nil { // token 语法错误
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": parseErr.Error(),
		})
		return
	}
	// 判断解析结果，验证通过 使用 .Valid 表示
	if !tokenObj.Valid { // 签名错误
		// 解析，验证失败
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Signature was error",
		})
		return
	}
	// token 验证通过

	// 获取 token 中包含的用户信息
	userName:=""
	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok {
		userName = claims["aud"].(string)
	}
	// 校验通过，获取用户名称，利用用户名称获取用户全部信息
	user := model.User{}
	orm.Where("user=?", userName).Find(&user)
	user.PasswordSalt = ""
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": user,
	})
}


func MemberLogin(c *gin.Context) {
	user := model.User{}
	// 先通过用户名获取用户信息，再比较密码是否正确
	postUser := c.PostForm("User")
	if orm.Where("user=?", postUser).Find(&user).RecordNotFound() {
		// 没有该用户
		c.JSON(http.StatusOK, gin.H{
			"error": "1 用户或密码错误，用户名",
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
			"error": "2 用户或密码错误，密码",
		})
		return
	}

	// 认证通过
	// 生成 JWT 格式的 Token
	// 生成签名的key
	mySigningKey := []byte(config.App["SECRET"])

	// token中包含如下的数据
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + (30 * 24 * 3600), // 有效期 当前时间戳+有效期
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
