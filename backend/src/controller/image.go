package controller

import (
	"config"
	"crypto/md5"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/gomango/imgtype"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"net/http"
	"os"
)

func ImageUpload(c *gin.Context) {

	file, fileErr := c.FormFile("file")
	if fileErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": fileErr.Error(),
		})
		return
	}
	//file.Header // 从浏览器的请求中获取的信息，是基于浏览器，是否可靠要取决于浏览器。

	f, _ := file.Open()
	defer f.Close()
	content := make([]byte, file.Size)
	f.Read(content)// 读取到文件的全部内容。

	// 检测类型
	mime, err := imgtype.DetectBytes(content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}




	// 确定了目标文件名
	//dst := fmt.Sprintf("%x", md5.Sum(content)) + filepath.Ext(file.Filename)
	dst := fmt.Sprintf("%x", md5.Sum(content))

	// 构建子目录
	subPath := string(dst[0]) + string(os.PathSeparator) + string(dst[1]) + string(os.PathSeparator)
	savePath := config.App["UPLOAD_PATH"] + subPath
	// 保证savePath 存在
	os.MkdirAll(savePath, 0755)

	// 为图像做重新编码，保证图像的正确性（安全性）。
	switch mime {
	case "image/jpeg":
		srcF, _ := file.Open()
		defer srcF.Close()
		img, err := jpeg.Decode(srcF)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": "1" + err.Error(),
			})
			return
		}
		dst += ".jpg"
		imgFile, err := os.Create(savePath + dst)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error":"2" +  err.Error(),
			})
			return
		}
		defer imgFile.Close()
		err = jpeg.Encode(imgFile, img, &jpeg.Options{Quality: 36})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": "3" + err.Error(),
			})
			return
		}

	case "image/png":
		srcF, _ := file.Open()
		defer srcF.Close()
		img, err := png.Decode(srcF)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": "1" + err.Error(),
			})
			return
		}
		dst += ".png"
		imgFile, err := os.Create(savePath + dst)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error":"2" +  err.Error(),
			})
			return
		}
		defer imgFile.Close()
		err = png.Encode(imgFile, img)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": "3" + err.Error(),
			})
			return
		}
	case "image/gif":
		srcF, _ := file.Open()
		defer srcF.Close()
		img, err := gif.Decode(srcF)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": "1" + err.Error(),
			})
			return
		}
		dst += ".gif"
		imgFile, err := os.Create(savePath + dst)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error":"2" +  err.Error(),
			})
			return
		}
		defer imgFile.Close()
		err = gif.Encode(imgFile, img, &gif.Options{NumColors: 256})
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error": "3" + err.Error(),
			})
			return
		}
	}

	// 制作缩略图
	small, err := MakeThumb(savePath + dst, 146, 146)
	if err != nil {

	}

	big, err := MakeThumb(savePath + dst, 1460, 1460)
	if err != nil {

	}

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": map[string]string{
			"Image": dst,
			"ImageSmall": small,
			"ImageBig": big,
		},
	})
	return
}

func MakeThumb(srcFile string, thumbW, thumbH int) (string, error) {
	// 打开原图像
	srcF, err := os.Open(srcFile)
	if err != nil {
		return "", err
	}


	// 完成计算
	// 获取图像的宽高
	srcConfig, _, err := image.DecodeConfig(srcF)
	srcF.Close()
	srcW, srcH := srcConfig.Width, srcConfig.Height

	var dstW, dstH int

	// 计算 宽之比 和 高之比
	if float64(srcW)/float64(thumbW) >= float64(srcH)/float64(thumbH) {
		dstW = thumbW // 宽与缩略图一致，
		dstH = int(math.Round(float64(dstW) * (float64(srcH)/float64(srcW)))) // 高等比计算出来
	} else {
		dstH = thumbH// 高与缩略图一致
		dstW = int(math.Round(float64(dstH) * (float64(srcW)/float64(srcH))))
	}

	// 计算位置
	var dstX, dstY int
	dstX = (thumbW - dstW) / 2
	dstY = (thumbH - dstH) / 2

	// 创建缩略图
	// 重新采样。1，从原图上采集新的点。2，重新编码为缩略图图片
	// 先利用一个成品包，将缩略图做出来，再放在我们的白色背景上

	// 创建缩略图
	thumbRect := image.Rect(0, 0, thumbW, thumbH)
	thumb := image.NewRGBA(thumbRect)
	//
	//// 使用白色进行背景填充
	bgColor := color.RGBA{
		255, 255, 255, 255,
	}
	// 填充，将白色画在画布上
	draw.Draw(thumb, thumbRect, &image.Uniform{C: bgColor}, image.Pt(0,0), draw.Src)
	//
	// 缩略图，将原图src图，画到thumb图上
	srcF1, err := os.Open(srcFile)
	if err != nil {
		return "", err
	}
	defer srcF1.Close()
	src, err := jpeg.Decode(srcF1) // 打开原图
	thumbSrc := imaging.Resize(src, dstW, dstH, imaging.Lanczos)

	srcRect := image.Rect(dstX, dstY, dstX+dstW, dstY+dstH)
	draw.Draw(thumb, srcRect, thumbSrc, image.Pt(0,0), draw.Src)

	tempFile := os.TempDir() + "/thumb.png"
	thumbFile, _ := os.Create(tempFile)
	png.Encode(thumbFile, thumb)
	thumbFile.Close()

	// 写入到指定文件
	// 利用文件内容形成摘要，存储到指定位置。
	content, _ := ioutil.ReadFile(tempFile)

	dst := fmt.Sprintf("%x", md5.Sum(content)) + ".png"
	// 构建子目录
	subPath := string(dst[0]) + string(os.PathSeparator) + string(dst[1]) + string(os.PathSeparator)
	savePath := config.App["UPLOAD_PATH"] + subPath
	// 保证savePath 存在
	os.MkdirAll(savePath, 0755)

	dstFile, err := os.Create(savePath + dst)
	dstFile.Write(content)
	dstFile.Close()

	return dst, nil
}
