/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-03 09:19:21
 * @FilePath: \image_hub\spiders\wechat_gh_8c96baecf453.go
 * @Description: 微信号自定义处理函数-头像有点好看
 */

package spiders

import (
	"image_hub/model"
	"strings"
)

// 微信号：gh_8c96baecf453, 文章索引: 1 的自定义处理函数
// 头像内容处理
func gh_8c96baecf453(sections []model.Section) []model.Section {

	// Filter out sections with invalid text
	sections = filterDirtyText(sections)

	// 1. 如果是头像, 壁纸, 朋友圈背景 过滤掉gif图
	// file://D:/work/wechat_download_data/html/Dump-0422-20-45-37/20230110_171808_1.html
	// file://D:/work/wechat_download_data/html/Dump-0422-20-45-37/20230105_180534_1.html
	// file://D:/work/wechat_download_data/html/Dump-0422-20-45-37/20230103_172345_2.html
	// file://D:/work/wechat_download_data/html/Dump-0422-20-45-37/20230109_150750_2.html

	// 2. 文本 "部位 女头" 去掉空格

	// 3. 如果后面文本是"女生头像号", 则删掉该section
	// file://D:/work/wechat_download_data/html/Dump-0422-20-45-37/20230122_222156_1.html

	for i := 0; i < len(sections); i++ {
		sections[i].Text = strings.Replace(sections[i].Text, " ", "", -1)
		if strings.Contains(sections[i].Text, "女生头像号") {
			sections = append(sections[:i], sections[i+1:]...)
			i-- // compensate for the removed element by decrementing the index
		}
	}

	// Filter out sections with empty image_urls
	sections = filterEmptyImageUrls(sections)

	return sections
}
