package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Require(pri string) gin.HandlerFunc {
	//pri
	return func(c *gin.Context) {
		// root 上帝视角
		userT, _ := c.Get("userName")
		userName := userT.(string)
		if ("root" == userName) {
			return
		}

		// 全部权限
		prisT,_ := c.Get("pris")
		pris := prisT.([]string)

		// 判断其中是否具有pri，需要的权限即可
		has := false
		for _, p := range pris {
			if p == pri {
				has = true
				break
			}
		}
		if !has {
			// 无权限
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"error": "privilege not found",
			})
			return
		}

	}
}
