package middleware

import (
	"bytes"
	"config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

// JWT 校验中间件
func JWTToken(c *gin.Context) {
	// 得到 TOKEN
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		// 没有 authorization 头
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "header Authorization not found",
		})
		c.Abort()
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
		c.Abort()
		return
	}
	// 判断解析结果，验证通过 使用 .Valid 表示
	if !tokenObj.Valid { // 签名错误
		// 解析，验证失败
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Signature was error",
		})
		c.Abort()
		return
	}
	// token 验证通过

	// 获取 token 中包含的用户信息
	userName:=""
	if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok {
		userName = claims["aud"].(string)
	}
	// 先在该中间件存储起来
	c.Set("userName", userName)

}