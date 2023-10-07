/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-10-06 18:14:26
 * @FilePath: \image_hub\spiders\wechat_seevanlove.go
 * @Description: 微信号自定义处理函数-情侣头像原创榜
 */

package spiders

import (
	"image_hub/model"
	"strings"
)

// 微信号：touxiangshe, 文章索引: 1 的自定义处理函数
// 头像内容处理
func seevanlove_1(_ *model.TblArticle, sections []model.Section) []model.Section {

	// Filter out sections with invalid text
	sections = filterDirtyText(sections)

	// 遍历头像当前面的section如果section.Text不为空且含有"头像", section.ImageUrls为空，则将该section.Text赋值为后面的section.ImageUrls不为空的section.Text
	// Concatenate sections containing "头像" with subsequent sections that have non-empty image_urls
	for i := 0; i < len(sections)-1; i++ {
		// Check if the current section has "头像" in its text and has an empty image_urls array
		if len(sections[i].ImageUrls) == 0 && strings.Contains(sections[i].Text, "头像") {
			// Look for the next section that has non-empty image_urls and concatenate the text with a newline separator
			for j := i + 1; j < len(sections); j++ {
				// Only concatenate sections with unique text and don't concatenate the same section more than once
				if sections[j].Text != sections[i].Text && len(sections[j].ImageUrls) > 0 {
					sections[j].Text = sections[i].Text
					break
				}
			}
		}
	}

	// Filter out sections with empty image_urls
	sections = filterEmptyImageUrls(sections)

	// 检查sections中，如果section.Text中有"爱 · 你 · 所 · 爱" 或 "爱你所爱" 文字时，则删掉该section
	for i := 0; i < len(sections); i++ {
		if strings.Contains(sections[i].Text, "爱 · 你 · 所 · 爱") || strings.Contains(sections[i].Text, "爱你所爱") {
			sections = append(sections[:i], sections[i+1:]...)
			i-- // compensate for the removed element by decrementing the index
		}
	}

	// 如果sections中sections[i]内的ImageUrls数量小于4时，将它与相邻的下一个sections[i+1]合并为一个，并继续检查sections[i+1]内的ImageUrls数量，直到ImageUrls数量大于或等于4。合并后，每个sections[i]内的ImageUrls数量都将大于或等于4
	sections = mergeImageUrls(sections)

	return sections
}
