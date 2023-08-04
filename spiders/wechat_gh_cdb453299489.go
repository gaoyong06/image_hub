/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-04 15:59:34
 * @FilePath: \image_hub\spiders\wechat_gh_cdb453299489.go
 * @Description: 微信号自定义处理函数-要啥头像
 */

package spiders

import (
	"image_hub/model"
	"strings"
)

// gh_cdb453299489, 文章索引: 1,2,3 的自定义处理函数
// 头像,壁纸,背景图, 广告内容处理
func gh_cdb453299489(article *model.TblArticle, sections []model.Section) []model.Section {

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

	// 基本是一个网页上一组内容,将所以的section合并到一个section
	var newSections []model.Section
	newSection := model.Section{
		Idx:       1,
		Text:      article.Title,
		ImageUrls: []string{},
	}

	for i := 0; i < len(sections); i++ {
		newSection.ImageUrls = append(newSection.ImageUrls, sections[i].ImageUrls...)
	}
	newSections = append(newSections, newSection)

	return newSections
}
