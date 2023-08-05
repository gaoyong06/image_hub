/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-05 07:48:35
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
	sections = filterOnlyOneImageUrls(sections)

	// Filter out sections with empty image_urls
	sections = filterEmptyImageUrls(sections)

	// 如果sections内image_urls都只有两个图片,将每4个sections内的item合并成一个, 删掉被合并的item， 合并后sections内item的image_urls都是8张图片
	sections = mergeImageUrls(sections)

	return sections
}
