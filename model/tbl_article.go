/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-24 11:15:14
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-26 14:20:50
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
	Mid         int       `json:"mid"`                             // 文章id 每篇文章的唯一标识符
	Biz         string    `json:"biz"`                             // 微信公众号的唯一标识符
	Idx         int       `json:"idx"`                             // 如果一篇文章有多页内容，idx表示当前页面是第几页
	Sn          string    `json:"sn"`                              // 一篇文章的唯一标识符，与mid不同的是，sn是加密后的标识符
	Title       string    `json:"title"`                           // 标题
	Author      string    `json:"author"`                          // 作者
	Tags        []string  `gorm:"serializer:json" json:"tags"`     // 合集标签
	Sections    []Section `gorm:"serializer:json" json:"sections"` // 文章分段，一篇文章(article)由多个分段(section)组成
	LocalPath   string    `json:"local_path"`                      // 文章保存路径
	PublishTime time.Time `json:"publish_time"`                    // 发布时间
}

func GetTblArticle() *TblArticle {
	return &TblArticle{}
}

func (t *TblArticle) TableName() string {
	return "tbl_article"
}

// https://gorm.io/zh_CN/docs/advanced_query.html
func (t *TblArticle) CreateOrUpdate() (int, error) {

	condModel := TblArticle{Mid: t.Mid}
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

	return t.Mid, err
}

// 另一种实现方式
// func (t *TblArticle) CreateOrUpdate() (int, error) {

// 	var oldArticle TblArticle
// 	var err error

// 	result := DB.Table(t.TableName()).Where("mid = ?", t.Mid).First(&oldArticle)
// 	if result.Error != nil {
// 		if result.Error == gorm.ErrRecordNotFound {
// 			err = DB.Table(t.TableName()).Create(t).Error
// 			return t.Mid, err
// 		}
// 		return t.Mid, result.Error
// 	}
// 	err = DB.Table(t.TableName()).Save(t).Error
// 	return t.Mid, err
// }
