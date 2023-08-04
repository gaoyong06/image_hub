/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-04 12:02:30
 * @FilePath: \image_hub\spiders\wechat_touxiangcool.go
 * @Description: 微信号自定义处理函数-你的小众头像
 */

package spiders

import (
	"image_hub/model"
	"strings"
)

// h13031h, 文章索引: 1,2,3,4 的自定义处理函数
// 头像,壁纸,背景图, 广告内容处理
func h13031h(_ *model.TblArticle, sections []model.Section) []model.Section {

	// Filter out sections with invalid text
	sections = filterDirtyText(sections)

	// 如果最后一个section的image_urls只有一个图的话则删掉
	// file://D:/work/wechat_download_data/html/Dump-0423-19-16-40/20200123_225558_1.html

	// 如果section的image_urls内只有一个图片的话,则删掉
	// file://D:/work/wechat_download_data/html/Dump-0423-19-16-40/20200130_120000_2.html

	// 删掉第一个section的内容
	// file://D:/work/wechat_download_data/html/Dump-0423-19-16-40/20200203_110653_1.html

	// 文本去除空格
	// file://D:/work/wechat_download_data/html/Dump-0423-19-16-40/20200422_103100_1.html

	for i := 0; i < len(sections); i++ {
		sections[i].Text = strings.Replace(sections[i].Text, " ", "", -1)
		if (sections[i].Text == "REC") || len(sections[i].ImageUrls) == 1 {
			sections = append(sections[:i], sections[i+1:]...)
			// compensate for the removed element by decrementing the index
			i--
		}
	}

	// Filter out sections with empty image_urls
	sections = filterEmptyImageUrls(sections)
	return sections
}
