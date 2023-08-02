/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-02 18:00:54
 * @FilePath: \image_hub\spiders\wechat_touxiangshe.go
 * @Description: 微信号自定义处理函数-头像社
 */

package spiders

import (
	"image_hub/model"
)

// 微信号：touxiangshe, 文章索引: 1 的自定义处理函数
// 头像内容处理
func touxiangshe_1(sections []model.Section) []model.Section {

	// 过滤字符串
	sections = filterDirtyText(sections)

	// 将最后一个ImageUrls不为空的Section中的Text赋值为紧邻他的ImageUrls为空的Text的值
	sections = modifyLastSectionText(sections)

	// 去掉sections中最后一个section中的imageUrls中的最后一项，然后返回修改后的sections
	if len(sections) > 0 {
		lastSection := &sections[len(sections)-1]
		if len(lastSection.ImageUrls) > 0 {
			lastSection.ImageUrls = lastSection.ImageUrls[:len(lastSection.ImageUrls)-1]
		}
	}

	// Filter out sections with empty image_urls
	sections = filterEmptyImageUrls(sections)

	return sections
}

// 微信号：touxiangshe, 文章索引: 2 的自定义处理函数
// 背景图内容处理
func touxiangshe_2(sections []model.Section) []model.Section {

	// 过滤字符串
	sections = filterDirtyText(sections)

	// 如果sections数组中第1项imageUrls中只有1个值，则将该项的imageUrls设置为空数组
	// 是一个大图
	if len(sections) > 1 && len(sections[0].ImageUrls) == 1 {
		sections[0].ImageUrls = []string{}
	}

	// 如果sections数组中第2项imageUrls为空数组，但Text不为空，则将该文本设置时sections中第3项的文本
	// 是大图下面的文案
	if len(sections) > 3 && len(sections[1].ImageUrls) == 0 && len(sections[1].Text) > 0 {
		sections[2].Text = sections[1].Text
	}

	// Filter out sections with empty image_urls
	sections = filterEmptyImageUrls(sections)

	return sections
}

// 微信号：touxiangshe, 文章索引: 3 的自定义处理函数
// 壁纸内容处理
func touxiangshe_3(sections []model.Section) []model.Section {

	// 过滤字符串
	sections = filterDirtyText(sections)

	// 将最后一个ImageUrls不为空的Section中的Text赋值为紧邻他的ImageUrls为空的Text的值
	sections = modifyLastSectionText(sections)

	// Filter out sections with empty image_urls
	sections = filterEmptyImageUrls(sections)

	return sections
}

// 微信号：touxiangshe, 文章索引: 4 的自定义处理函数
// 表情包内容处理
func touxiangshe_4(sections []model.Section) []model.Section {

	// 过滤字符串
	sections = filterDirtyText(sections)

	// 如果sections数组中第1项imageUrls中只有1个值，则将该项的imageUrls设置为空数组
	// 是一个emoji表情
	if len(sections) > 1 && len(sections[0].ImageUrls) == 1 {
		sections[0].ImageUrls = []string{}
	}

	// Filter out sections with empty image_urls
	sections = filterEmptyImageUrls(sections)

	return sections
}

//------------------------------- 公共方法 ---------------------------------------------

// 将最后一个ImageUrls不为空的Section中的Text赋值为紧邻他的ImageUrls为空的Text的值
func modifyLastSectionText(sections []model.Section) []model.Section {

	// 将最后一个ImageUrls不为空的Section中的Text赋值为紧邻他的ImageUrls为空的Text的值
	// 查找最后一个ImageUrls不为空的Section
	lastImageSectionIndex := -1
	for i := len(sections) - 1; i >= 0; i-- {
		if len(sections[i].ImageUrls) > 0 {
			lastImageSectionIndex = i
			break
		}
	}

	// 如果找到了最后一个ImageUrls不为空的Section
	if lastImageSectionIndex != -1 {

		// 查找最后一个ImageUrls为空的Section
		lastEmptySectionIndex := -1
		for i := lastImageSectionIndex + 1; i < len(sections); i++ {
			if len(sections[i].ImageUrls) == 0 && len(sections[i].Text) > 0 {
				lastEmptySectionIndex = i
				break
			}
		}

		// 如果找到了最后一个ImageUrls为空的Section
		if lastEmptySectionIndex != -1 {

			// 将最后一个ImageUrls不为空的Section中的Text赋值为紧邻他的ImageUrls为空的Text的值
			sections[lastImageSectionIndex].Text = sections[lastEmptySectionIndex].Text
		}
	}

	return sections
}
