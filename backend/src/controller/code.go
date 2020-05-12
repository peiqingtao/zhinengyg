package controller

import (
	"crypto/md5"
	"fmt"
	"strings"

	//"crypto/md5"
	"github.com/afocus/captcha"
	"github.com/gin-gonic/gin"
	"image/png"
	_ "image/png"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

func CheckCode(c *gin.Context) {
	// 设置 captcha 对象
	cap := captcha.New()
	font := "./src/config/BRADHITC.TTF"
	cap.SetFont(font)

	// 4 个长度的数值验证码
	img, code:= cap.Create(4,captcha.NUM)
	// 将图像作为响应主体
	c.Header("Content-Type", "image/png")
	png.Encode(c.Writer, img)

	// 将验证码码值存储起来（为了验证）
	//127.0.0.1:64756 + [Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:69.0) Gecko/20100101 Firefox/69.0
	i := strings.Index(c.Request.RemoteAddr, "]")
	key := fmt.Sprintf("%x", md5.Sum([]byte(c.Request.RemoteAddr[1:i] + c.Request.Header["User-Agent"][0])))
	// 临时存储，建议使用缓存系统，因为验证码通常具有有效期，例如2分钟
	Rds.Do("SET", "code_" + key, code)
	Rds.Do("expire", "code_"+key, 2 * 60)
}

func SmsCode(c *gin.Context) {
	// 获取短信发送的客户端，提供设置的
	accessKeyId := "LTAI4FmR9Gcj5YDoj1uBYTcE"
	accessSecret := "BK370QcLmd1zyGs7zFgvsYgRJF7BAM"
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessSecret)

	// 配置请求
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	tel := "15101542640"
	request.PhoneNumbers = tel
	signName := "Firelinks"
	request.SignName = signName
	templateCode := "SMS_175533096"
	request.TemplateCode = templateCode
	code := "9527"
	request.TemplateParam = fmt.Sprintf(`{"code":"%s"}`, code)

	response, err := client.SendSms(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)

}
