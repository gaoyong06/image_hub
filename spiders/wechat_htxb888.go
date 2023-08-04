/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-04 15:27:53
 * @FilePath: \image_hub\spiders\wechat_htxb888.go
 * @Description: 微信号自定义处理函数-你的小众头像
 */

package spiders

import (
	"image_hub/model"
	"strings"
)

// htxb888, 文章索引: 1,2,3,4 的自定义处理函数
// 头像,壁纸,背景图, 广告内容处理
func htxb888(_ *model.TblArticle, sections []model.Section) []model.Section {

	// Filter out sections with invalid text
	sections = filterDirtyText(sections)

	// 文本去除空格
	// 如果section的image_urls内只有一个图片的话,则删掉

	for i := 0; i < len(sections); i++ {
		sections[i].Text = strings.Replace(sections[i].Text, " ", "", -1)
		if len(sections[i].ImageUrls) == 1 {
			sections = append(sections[:i], sections[i+1:]...)
			// compensate for the removed element by decrementing the index
			i--
		}
	}

	// Filter out sections with empty image_urls
	sections = filterEmptyImageUrls(sections)

	return sections
}