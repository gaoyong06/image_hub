/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-03 17:42:07
 * @FilePath: \image_hub\spiders\wechat_touxiangcool.go
 * @Description: 微信号自定义处理函数-头像库
 */

package spiders

import (
	"image_hub/model"
	"strings"
)

// touxiangcool, 文章索引: 1,2,3,4 的自定义处理函数
// 头像,壁纸,背景图, 广告内容处理
func touxiangcool(_ *model.TblArticle, sections []model.Section) []model.Section {

	// Filter out sections with invalid text
	sections = filterDirtyText(sections)

	// 删掉壁纸第一个gif
	// file://D:/work/wechat_download_data/html/Dump-0423-11-39-39/20220802_210000_3.html

	// 第一个和后面的内容合并
	// file://D:/work/wechat_download_data/html/Dump-0423-11-39-39/20220801_210000_4.html

	// Filter out sections with empty image_urls
	sections = filterEmptyImageUrls(sections)

	return sections
}

func touxiangcool_4(article *model.TblArticle, sections []model.Section) []model.Section {

	// Filter out sections with invalid text
	sections = filterDirtyText(sections)

	// 删掉壁纸第一个gif
	// file://D:/work/wechat_download_data/html/Dump-0423-11-39-39/20220802_210000_3.html

	// 第一个和后面的内容合并
	// file://D:/work/wechat_download_data/html/Dump-0423-11-39-39/20220801_210000_4.html

	var texts []string
	texts = append(texts, article.Title)
	for i := 0; i < len(sections); i++ {

		// 去掉第一个图片内容
		if len(sections[i].ImageUrls) == 1 {
			sections[i].ImageUrls = []string{}
		}

		// 将文章标题和section内多个图片之前的所有文案，合并在一起
		texts = append(texts, sections[i].Text)
		if len(sections[i].ImageUrls) > 1 {
			sections[i].Text = strings.Join(texts, "\n")
			break
		}
	}

	// Filter out sections with empty image_urls
	sections = filterEmptyImageUrls(sections)
	return sections
}
