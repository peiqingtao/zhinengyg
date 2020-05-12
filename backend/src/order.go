package main

import (
	"config"
	"controller"
	"encoding/json"
	"log"
	"model"
	"time"
)

func main() {
	log.Println("订单处理中...")

	config.InitConfig()

	db, err := controller.InitDB()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	// redis
	rc, err := controller.InitRedis()
	if err != nil {
		log.Println(err)
		return
	}
	defer rc.Close()

	for {
		// 从消息队列，获取订单信息，进行处理。
		result, err := controller.Rds.Do("XREAD", "COUNT", "1", "BLOCK", "30000", "STREAMS", "orderQueue", "$")
		if err != nil {
			log.Println(err)
			continue
		}
		if result == nil {
			// 当前没有新的订单
			// 继续等待
			continue
		}

		// 解析result，获取订单sn
		sn := string(result.([]interface{})[0].([]interface{})[1].([]interface{})[0].([]interface{})[1].([]interface{})[1].([]byte))
		log.Println(sn)
		// 获取订单信息
		orderResult, err := controller.Rds.Do("HGET", "tempOrder", sn)
		if err != nil {
			log.Println(err)
			continue
		}
		orderInfo := controller.TempOrder{}
		parseErr := json.Unmarshal(orderResult.([]byte), &orderInfo)
		if parseErr != nil {
			log.Println(parseErr)
			continue
		}

		// 订单库存检测
		// 获取购物车中产品的购买量
		cart := model.Cart{}
		db.Where("user_id=?", orderInfo.UserID).Find(&cart)
		cartInfo := []struct{
			ProductID uint  `json:"productID"`
			BuyQuantity int `json:"buyQuantity"`
		}{}
		json.Unmarshal([]byte(cart.Content), &cartInfo)// 全部的购买商品【1, 2, 3, 4, 6】

		// 获取购买数量， 仅仅需要购买的部分 [1, 2]
		buyInfo  := []struct{
			ProductID uint  `json:"productID"`
			BuyQuantity int `json:"buyQuantity"`
		}{}
		for _, p := range cartInfo {
			for _, pID := range orderInfo.BuyProductID {
				if p.ProductID == pID {
					buyInfo = append(buyInfo, p)
					break
				}
			}
		}

		// 检测所需要购买的产品是否在购物车中
		if len(buyInfo) != len(orderInfo.BuyProductID) {
			log.Println("所购买的产品不在购物车中")
			controller.Rds.Do("HSET", "orderResult", sn, "error")
			continue
		}

		// 检测库存
		flag := true // 检测通过
		for _, p := range buyInfo {
			// 读取每个产品的库存，进行检测
			quantity := model.Quantity{}
			db.Where("product_id=?", p.ProductID).Find(&quantity)
			if quantity.Number < p.BuyQuantity {
				// 系统库存，小于购买库存，则失败
				flag = false // 未通过
				break
			}
		}
		if ! flag {
			// 库存不足，订单失败
			// 记录结果
			log.Println("订单失败，库存不足")
			controller.Rds.Do("HSET", "orderResult", sn, "error")
			// 记录结果
			continue
		}

		// 库存充足
		log.Println("订单成功，库存充足")

		// 扣减库存，形成订单数据
		for _, p := range buyInfo {
			// 读取每个产品的库存，进行检测
			quantity := model.Quantity{}
			db.Where("product_id=?", p.ProductID).Find(&quantity)
			// 扣减库存
			quantity.Number = quantity.Number - p.BuyQuantity
			db.Save(&quantity)
		}

		// 形成订单数据
		order := model.Order{}
		order.Sn = sn
		order.AddressID = orderInfo.AddressID
		order.PaymentStatusID = 2 // 未支付
		order.ShippingID = orderInfo.ShippingId
		order.ShippingStatusID = 2 // 未配货
		order.Amount = 99 //
		order.OrderTime = time.Now()
		order.OrderStatusID = 2// 确认
		db.Create(&order)

		// 记录订单中的产品
		for _, p := range buyInfo {
			// 读取每个产品的库存，进行检测
			// 查询产品当前的价格
			product := model.Product{}
			db.Where("id=?", p.ProductID).Find(&product)

			orderProduct := model.OrderProduct{}
			orderProduct.OrderID = order.ID
			orderProduct.ProductID = p.ProductID
			orderProduct.BuyQuantity = p.BuyQuantity
			orderProduct.BuyPrice = int(product.Price * 100)
			db.Create(&orderProduct)
		}

		// 删除购物车中已购产品
		restInfo := []struct{
			ProductID uint  `json:"productID"`
			BuyQuantity int `json:"buyQuantity"`
		}{}
		for _, p := range cartInfo {
			flag := true // 假定为需要保留的
			for _, pID := range orderInfo.BuyProductID {
				if pID == p.ProductID {
					flag = false // 该产品是已经购买的
					break
				}
			}
			if flag {
				restInfo = append(restInfo, p)
			}
		}
		// 更新该购物车字段即可
		c, _ := json.Marshal(restInfo)
		cart.Content = string(c)
		db.Save(&cart)

		// 记录成功结果
		controller.Rds.Do("HSET", "orderResult", sn, "success")
	}
}
//[
//	[
//		[111 114 100 101 114 81 117 101 117 101]
//		[
//			[
//				[49 53 55 49 50 56 52 57 53 54 57 48 56 45 48]
//				[
//					[99 111 110 116 101 110 116]
//					[50 48 49 57 49 48 49 55 48 48 48 48 48 48 48 54]
//				]
//			]
//		]
//	]
//]
