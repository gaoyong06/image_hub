/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-24 11:15:14
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-24 11:15:14
 * @FilePath: \image_hub\model\section.go
 * @Description: 公众号分段信息
 */
package model

// Section 即一篇公众号文章内的一个分段
// 一个分段是：标题(或 段落)文字 + 下方的图片列表
type Section struct {
	Idx       int      `json:"idx"`                               // 分段索引
	Text      string   `json:"text"`                              // 分段内文字内容
	ImageUrls []string `gorm:"serializer:json" json:"image_urls"` // 分段内图片url列表
}
