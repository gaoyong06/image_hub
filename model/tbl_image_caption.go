package model

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// TblImageCaption 图片描述信息表
type TblImageCaption struct {
	ImageId     int64     `json:"image_id" gorm:"image_id"`         // 唯一id
	LocalPath   string    `json:"local_path" gorm:"local_path"`     // 图片本地保存路径
	Width       int64     `json:"width" gorm:"width"`               // 图片宽度
	Height      int64     `json:"height" gorm:"height"`             // 图片高度
	AspectRatio string    `json:"aspect_ratio" gorm:"aspect_ratio"` // 图片宽高比
	Format      string    `json:"format" gorm:"format"`             // 图片格式
	Size        float64   `json:"size" gorm:"size"`                 // 图片大小
	Md5         string    `json:"md5" gorm:"md5"`                   // md5
	Phash       string    `json:"phash" gorm:"phash"`               // phash
	OcrText     string    `json:"ocr_text" gorm:"ocr_text"`         // ocr识别文本 null表示未识别, 空字符串和文本表示识别的结果
	Md5Count    int64     `json:"md5_count" gorm:"md5_count"`       // 相同md5的图片数量
	PhashCount  int64     `json:"phash_count" gorm:"phash_count"`   // 相同phash的图片数量
	Tags        string    `json:"tags" gorm:"tags"`                 // 图片识别的图片标签
	Caption     string    `json:"caption" gorm:"caption"`           // 图片识别的内容描述
	IsAd        int8      `json:"is_ad" gorm:"is_ad"`               // 是否是广告图 初始值：-1, 否：0, 是：1
	AdReason    string    `json:"ad_reason" gorm:"ad_reason"`       // 图片被定义为广告的原因
	IsQrcode    int8      `json:"is_qrcode" gorm:"is_qrcode"`       // 是否包含二维码 初始值：-1, 否：0, 是：1
	IsDataset   int8      `json:"is_dataset" gorm:"is_dataset"`     // 是否是weaviate的dataset 初始值：-1, 否：0, 是：1
	CreatedAt   time.Time `json:"created_at" gorm:"created_at"`     // 创建时间
	UpdatedAt   time.Time `json:"updated_at" gorm:"updated_at"`     // 最后修改时间
	DeletedAt   time.Time `json:"deleted_at" gorm:"deleted_at"`     // 删除时间
}

func GetTblImageCaption() *TblImageCaption {
	return &TblImageCaption{}
}

// TableName 表名称
func (t *TblImageCaption) TableName() string {
	return "tbl_image_caption"
}

// 判断图片是否合法
func (t *TblImageCaption) IsValidImage(localPath string) bool {
	// 根据 local_path 查询对应的记录
	if err := contentNerDB.Where("local_path = ?", localPath).First(t).Error; err != nil {
		log.Error("Failed to query TblImageCaption:", err)
		return false
	}

	// 判断 is_ad 是否等于 1 或者 is_qrcode 是否等于 1
	// 非法图片
	if t.IsAd == 1 || t.IsQrcode == 1 {
		return false
	}

	// 合法图片
	return true
}
