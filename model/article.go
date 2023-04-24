/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-24 11:15:14
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-24 11:16:07
 * @FilePath: \image_hub\model\article.go
 * @Description: 公众号文章信息
 */
package model

import "time"

// Article 即一篇公众号文章内容
type Article struct {
	ArticleId   int       `json:"article_id"`   // 文章id
	Title       string    `json:"title"`        // 标题
	Author      string    `json:"author"`       // 作者
	Tags        []string  `json:"tag"`          // 合集标签
	Sections    []Section `json:"sections"`     // 文章分段，一篇文章(article)由多个分段(section)组成
	PublishTime time.Time `json:"publish_time"` // 发布时间
}
