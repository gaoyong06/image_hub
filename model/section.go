/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-03-09 21:54:28
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-21 17:18:57
 * @FilePath: \image_hub\model\car_info.go
 * @Description: 二手车信息
 */
package model

// Section 即一篇公众号文章内的一个分段
// 一个分段是：标题(或 段落)文字 + 下方的图片列表
type Section struct {
	SectionId int    `json:"section_id"` // 分段id
	Text      string `json:"text"`       // 分段内文字内容
	ImageUrls string `json:"image_urls"` // 分段内图片url列表
}
