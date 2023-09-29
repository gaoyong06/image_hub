/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-24 11:15:14
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-09-29 11:23:25
 * @FilePath: \image_hub\model\article.go
 * @Description: 公众号文章信息
 */
package model

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Article 即一篇公众号文章内容
type TblArticle struct {
	Sn          string    `gorm:"type:VARCHAR(255);column:sn;primary_key USING BTREE" json:"sn"`                                                     // 一篇文章的唯一标识符
	Mid         uint64    `gorm:"type:BIGINT;column:mid;default:0;NOT NULL" json:"mid"`                                                              // 每次推送文章的唯一标识符
	Idx         uint      `gorm:"type:INT(11);column:idx;default:0;NOT NULL" json:"idx"`                                                             // 如果一次推送有多篇文章，idx表示当前页面是第几个
	Biz         string    `gorm:"type:VARCHAR(255);column:biz;NOT NULL" json:"biz"`                                                                  // 微信公众号的唯一标识符
	Author      string    `gorm:"type:VARCHAR(255);column:author;NOT NULL" json:"author"`                                                            // 公众号作者名称
	WechatId    string    `gorm:"type:VARCHAR(255);column:wechat_id;NOT NULL" json:"wechat_id"`                                                      // 公众号微信号
	Title       string    `gorm:"type:VARCHAR(255);column:title;NOT NULL" json:"title"`                                                              // 文章标题
	Tags        []string  `gorm:"type:MEDIUMTEXT;column:tags;serializer:json" json:"tags"`                                                           // 合集标签
	Sections    []Section `gorm:"type:MEDIUMTEXT;column:sections;serializer:json" json:"sections"`                                                   // 文章分段，一篇文章(article)由多个分段(section)组成
	LocalPath   string    `gorm:"type:VARCHAR(255);column:local_path;NOT NULL" json:"local_path"`                                                    // 文章本地保存路径
	PublishTime time.Time `gorm:"type:TIMESTAMP;column:publish_time;default:'0000-00-00 00:00:00';NOT NULL" json:"publish_time"`                     // 文章发布时间
	CreatedAt   time.Time `gorm:"type:TIMESTAMP;column:created_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"created_at"`                             // 创建时间
	UpdatedAt   time.Time `gorm:"type:TIMESTAMP;column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;NOT NULL" json:"updated_at"` // 最后修改时间
	DeletedAt   time.Time `gorm:"type:TIMESTAMP;column:deleted_at;default:'0000-00-00 00:00:00';NOT NULL" json:"deleted_at"`                         // 删除时间
}

func GetTblArticle() *TblArticle {
	return &TblArticle{}
}

func (t *TblArticle) TableName() string {

	// 表名为: tbl_article_微信号
	// 注意: tbl_article_后缀名 中,后缀名为微信号，但是如果微信号中如果包含中划线"-", 会将中划线"-"替换为下划线"_"
	wechatId := strings.Replace(t.WechatId, "-", "_", -1)

	tableName := fmt.Sprintf("%s_%s", "tbl_article", wechatId)
	return tableName
}

func (t *TblArticle) CreateTableIfNotExists() error {

	tableName := t.TableName()

	// Check if the table already exists
	if DB.Migrator().HasTable(tableName) {
		return nil
	}

	// Create the table
	err := DB.Set("gorm:table_options", "ENGINE=InnoDB").Table(tableName).AutoMigrate(t)
	if err != nil {
		log.Errorf("Failed to create table %s. Error: %+v\n", tableName, err)
		return err
	}

	return nil
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
		WechatId:    t.WechatId,
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
