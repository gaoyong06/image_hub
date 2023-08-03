/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-03 17:15:33
 * @FilePath: \image_hub\spiders\wechat_gh_22c17e1db325.go
 * @Description: 微信号自定义处理函数-头像即新欢
 */

package spiders

import (
	"fmt"
	"image_hub/model"
	"strings"
)

// gh_22c17e1db325, 文章索引: 1,2,3,4 的自定义处理函数
// 头像,壁纸,背景图, 广告内容处理
func gh_22c17e1db325(_ *model.TblArticle, sections []model.Section) []model.Section {

	// Filter out sections with invalid text
	sections = filterDirtyText(sections)

	// 广告处理
	// file:///D:/work/wechat_download_data/html/Dump-0422-20-54-12/20230418_135629_1.html
	// file:///D:/work/wechat_download_data/html/Dump-0422-20-54-12/20230318_131131_2.html

	for i := 0; i < len(sections)-1; i++ {

		// 遍历头像当前面的section如果section.Text不为空且含有"往下滑有"和"头像",则将该section.Text设置为空
		if strings.Contains(sections[i].Text, "往下滑有") && strings.Contains(sections[i].Text, "头像") {
			sections[i].Text = ""
		}
	}

	// 遍历头像当前面的section如果section.Text以"头像"结尾，则将该section.Text追加到后面的section.ImageUrls不为空的section.Text
	for i := 0; i < len(sections)-1; i++ {

		if strings.HasSuffix(sections[i].Text, "头像") {
			for j := i + 1; j < len(sections); j++ {

				// 如果text文本以"头像"结尾,则跳出循环
				if strings.HasSuffix(sections[j].Text, "头像") {
					break
				}

				// 拼接前后的text文本字符串
				if sections[j].Text != sections[i].Text && len(sections[j].ImageUrls) > 0 && !strings.Contains(sections[j].Text, "头像") {
					sections[j].Text = fmt.Sprintf("%s|%s", sections[i].Text, sections[j].Text)
				}

				// 删掉文本后面的|
				sections[j].Text = strings.TrimSuffix(sections[j].Text, "|")
			}
		}
	}

	// 如果第一个(或其他)section的image_urls，只有一个图片，则删掉该section
	for i := 0; i < len(sections)-1; i++ {

		if len(sections[i].ImageUrls) == 1 {
			sections[i].ImageUrls = []string{}
		}
	}

	// Filter out sections with empty image_urls
	sections = filterEmptyImageUrls(sections)

	return sections
}
