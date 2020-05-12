package model

import "github.com/jinzhu/gorm"

// 产品图像 模型
type Image struct {
	gorm.Model

	ProductID uint
	IsDefault bool
	SortOrder int
	Host string
	Image string // 从客户端上传后，保存到服务器的字段
	ImageSmall string // 图像缩略图路径，146 * 146
	ImageBig string // 图像缩略图路径， 1460* 1460
}

