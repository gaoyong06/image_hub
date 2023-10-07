/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-10-06 18:55:59
 * @FilePath: \image_hub\spiders\wechat_default.go
 * @Description: 微信号自定义处理函数-默认微信账号处理程序
 */

package spiders

import (
	"image_hub/model"
)

// default, 文章索引: 1的自定义处理函数
// 头像,壁纸,背景图, 广告内容处理
func simple(_ *model.TblArticle, sections []model.Section) []model.Section {

	// Filter out sections with invalid text
	sections = filterDirtyText(sections)

	// 文本去除空格
	sections = replaceTextBlank(sections)

	// 如果section的image_urls内只有一个图片的话,则删掉
	// sections = filterOnlyOneImageUrls(sections)

	// Filter out sections with empty image_urls
	sections = filterEmptyImageUrls(sections)

	// 如果sections中sections[i]内的ImageUrls数量小于4时，将它与相邻的下一个sections[i+1]合并为一个，并继续检查sections[i+1]内的ImageUrls数量，直到ImageUrls数量大于或等于4。合并后，每个sections[i]内的ImageUrls数量都将大于或等于4
	sections = mergeImageUrls(sections)

	return sections
}
