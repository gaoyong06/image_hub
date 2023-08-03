/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-24 11:15:14
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-03 15:10:28
 * @FilePath: \image_hub\model\article.go
 * @Description: 公众号文章信息
 */
package model

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Article 即一篇公众号文章内容
type TblArticle struct {
	Sn          string    `json:"sn"`                              // 一篇文章的唯一标识符
	Mid         int       `json:"mid"`                             // 每次推送文章的唯一标识符
	Idx         int       `json:"idx"`                             // 如果一次推送有多篇文章，idx表示当前页面是第几个
	Biz         string    `json:"biz"`                             // 微信公众号的唯一标识符
	Author      string    `json:"author"`                          // 公众号作者名称
	Title       string    `json:"title"`                           // 文章标题
	Tags        []string  `gorm:"serializer:json" json:"tags"`     // 合集标签
	Sections    []Section `gorm:"serializer:json" json:"sections"` // 文章分段，一篇文章(article)由多个分段(section)组成
	LocalPath   string    `json:"local_path"`                      // 文章本地保存路径
	PublishTime time.Time `json:"publish_time"`                    // 文章发布时间
}

func GetTblArticle() *TblArticle {
	return &TblArticle{}
}

func (t *TblArticle) TableName() string {
	return "tbl_article_touxiangcool"
}

// https://gorm.io/zh_CN/docs/advanced_query.html
func (t *TblArticle) CreateOrUpdate() (string, error) {

	condModel := TblArticle{Sn: t.Sn}
	assignModel := TblArticle{
		Biz:         t.Biz,
		Idx:         t.Idx,
		Sn:          t.Sn,
		Title:       t.Title,
		Author:      t.Author,
		Tags:        t.Tags,
		Sections:    t.Sections,
		LocalPath:   t.LocalPath,
		PublishTime: t.PublishTime,
	}

	err := DB.Table(t.TableName()).Where(condModel).Assign(assignModel).FirstOrCreate(t).Error
	if err != nil {
		log.Errorf("TblArticle CreateOrUpdate failed. err: %+v\n", err.Error())
	}

	return t.Sn, err
}

// 另一种实现方式
// func (t *TblArticle) CreateOrUpdate() (string, error) {

// 	var oldArticle TblArticle
// 	var err error

// 	result := DB.Table(t.TableName()).Where("sn = ?", t.Sn).First(&oldArticle)
// 	if result.Error != nil {
// 		if result.Error == gorm.ErrRecordNotFound {
// 			err = DB.Table(t.TableName()).Create(t).Error
// 			return t.Sn, err
// 		}
// 		return t.Sn, result.Error
// 	}
// 	err = DB.Table(t.TableName()).Save(t).Error
// 	return t.Sn, err
// }
