package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TempOrder struct {
	AddressID uint
	BuyProductID []uint
	ShippingId uint
	UserID uint
}


func OrderResult(c *gin.Context) {
	sn := c.DefaultQuery("sn", "")
	result, err := Rds.Do("HGET", "orderResult", sn)
	if err != nil || result == nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "order is not exists",
		})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": string(result.([]byte)),
	})

}

//订单创建
func OrderCreate(c *gin.Context) {

	user := member(c)
	if user == nil {
		// 未认证通过
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// 1 生成订单SN
	// 采用 日期+当天的序号的方案。
	//time.Time{}.Format()
	// 得到redis中的计数器的key，使用天作为标识
	now := time.Now()
	key := fmt.Sprintf("%d%d%d",
		now.Year(),
		now.Month(),
		now.Day(),
	)
	// 在redis使用的序号递增key
	counterKey := "counter" + key
	n, err := Rds.Do("incr", counterKey)
	if err != nil {
		// 没有序号，随机序号
		n = rand.Int63n(100000000)
	}
	// 变为字符串类型
	ns := strconv.Itoa(int(n.(int64)))
	// 固定长度补齐即可
	if l:=len(ns); l < 8 {
		ns = strings.Repeat("0", 8-l) + ns
	}
	// 订单号
	sn := key + ns

	// 2 临时存储订单信息
	tempOrder := TempOrder{}
	c.ShouldBind(&tempOrder)
	// 记录用户信息
	tempOrder.UserID = user.ID
	//记录产品信息

	// 将其json字符串话
	toj ,err := json.Marshal(tempOrder)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "提交的订单数据错误",
		})
		return
	}
	// 利用 hash 结构，存储在redis中。sn 为key，toj 为值
	_, hseterr := Rds.Do("HSET", "tempOrder", sn, toj)
	if hseterr != nil {
		// redis 缓存存储失败，应该使用其他的方案
	}

	// 3 将订单放入队列（等待处理）
	// 先获取队列长度，提示
	waitLen, _ := Rds.Do("XLEN", "orderQueue")

	//放入队列
	_, addErr := Rds.Do("XADD", "orderQueue", "*", "content", sn)
	if addErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "订单队列错误" + addErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": sn,
		"waitLen": waitLen,
	})
}

// 配送方式列表
func ShippingList(c *gin.Context) {
	user := member(c)
	if user == nil {
		// 未认证通过
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	data := []model.Shipping{}
	orm.Where("status=1").Find(&data)

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": data,
	})
}
