/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-24 11:15:14
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-26 11:51:10
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
	Biz         string    `json:"__biz"`        //  微信公众号的唯一标识符
	Mid         int       `json:"mid"`          // 文章id 每篇文章的唯一标识符
	Idx         int       `json:"idx"`          // 如果一篇文章有多页内容，idx表示当前页面是第几页
	Sn          string    `json:"sn"`           // 标题
	Title       string    `json:"title"`        // 标题
	Author      string    `json:"author"`       // 作者
	Tags        []string  `json:"tag"`          // 合集标签
	Sections    []Section `json:"sections"`     // 文章分段，一篇文章(article)由多个分段(section)组成
	LocalPath   string    `json:"local_path"`   // 文章保存路径
	PublishTime time.Time `json:"publish_time"` // 发布时间
}

func GetTblArticle() *TblArticle {
	return &TblArticle{}
}

func (t *TblArticle) TableName() string {
	return "tbl_seller_phone"
}

func (t *TblArticle) CreateOrUpdate() (string, error) {

	// https://gorm.io/zh_CN/docs/advanced_query.html
	condModel := TblArticle{Sn: t.Sn}
	assignModel := TblArticle{

		Biz:         t.Biz,
		Mid:         t.Mid,
		Idx:         t.Idx,
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
