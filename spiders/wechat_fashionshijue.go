/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-03 20:27:11
 * @FilePath: \image_hub\spiders\wechat_fashionshijue.go
 * @Description: 微信号自定义处理函数-头像文案
 */

package spiders

import (
	"image_hub/model"
)

// fashionshijue, 文章索引: 1,2,3,4 的自定义处理函数
// 头像,壁纸,背景图, 广告内容处理
func fashionshijue(_ *model.TblArticle, sections []model.Section) []model.Section {

	// Filter out sections with invalid text
	sections = filterDirtyText(sections)

	// Filter out sections with empty image_urls
	sections = filterEmptyImageUrls(sections)

	return sections
}
