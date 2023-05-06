/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-28 17:33:09
 * @FilePath: \image_hub\spiders\func_map.go
 * @Description: 爬虫相关公用方法
 */

package spiders

import (
	"image_hub/model"
	"log"
	"strings"

	"golang.org/x/net/html"
)

// 从HTML字符串中解析出Section数组，包含文字和图片
func ParseSectionsFromHTML(htmlStr string) []model.Section {

	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		log.Fatal(err)
	}

	var sections []model.Section

	// 字符串过滤器，过滤掉不需要的标签，包括空的 span、不可见文本元素等, #activity-name，#meta_content，#js_tags 三个标签的过滤
	filter := func(n *html.Node) bool {

		if n.Type == html.ElementNode && n.Data == "script" {
			return true
		}
		if n.Type == html.ElementNode && n.Data == "style" {
			return true
		}
		if n.Type == html.ElementNode && n.Data == "head" {
			return true
		}
		if n.Type == html.ElementNode && n.Data == "title" {
			return true
		}
		if n.Type == html.ElementNode && n.Data == "meta" {
			return true
		}

		if n.Type == html.ElementNode && len(n.Attr) > 0 {
			for _, attr := range n.Attr {
				if attr.Key == "id" && (attr.Val == "activity-name" || attr.Val == "meta_content" || attr.Val == "js_tags") {
					return true
				}
			}
		}

		if n.Type == html.TextNode && strings.TrimSpace(n.Data) == "\u200d" {
			return true
		}
		return false
	}

	var parseNode func(*html.Node, bool)
	parseNode = func(n *html.Node, skip bool) {
		if filter(n) {
			skip = true
		} else if skip {
			return
		} else if n.Type == html.ElementNode && n.Data == "img" {

			// 如果当前节点为img标签，提取其中的src属性作为Section的图片url
			var imageUrl string
			for _, attr := range n.Attr {

				if attr.Key == "src" {

					imageUrl = attr.Val
				}
			}

			// 将图片url添加到当前Section的ImageUrls列表
			if len(sections) <= 0 {
				sections = append(sections, model.Section{
					Text:      "",
					ImageUrls: []string{},
				})
			}

			// imageUrl不为空则追加到ImageUrls中
			if len(imageUrl) > 0 {
				currentSection := sections[len(sections)-1]
				currentSection.ImageUrls = append(currentSection.ImageUrls, imageUrl)
				sections[len(sections)-1] = currentSection
			}

		} else if n.Type == html.TextNode && strings.TrimSpace(n.Data) != "" && strings.TrimSpace(n.Data) != "\u200d" {

			// 如果当前节点为文本节点，提取其中的文字内容作为Section的文本内容
			currentText := strings.TrimSpace(n.Data)

			// 创建一个新的Section，并添加到数组中
			newSection := model.Section{
				Text:      currentText,
				ImageUrls: []string{},
			}
			sections = append(sections, newSection)
		}

		// 递归调用parseNode处理当前节点的所有子节点
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNode(c, skip)
		}
	}

	// 从根节点开始遍历
	parseNode(doc, false)
	return sections
}

//------------------------------------ 私有方法 -------------------------------------------------

// 获取文件名最后的数字
// file://D:/work/wechat_download_data/html/test4/20220810_111900_1.html
func getFileName(filePath string) string {

	// 将文件路径按照"/"分割成数组
	arr := strings.Split(filePath, "/")
	// 获取数组最后一个元素
	last := arr[len(arr)-1]
	// 将最后一个元素按照"."分割成数组
	arr2 := strings.Split(last, ".")
	// 获取数组第一个元素
	fileName := arr2[0]
	// 将文件名最后的数字提取出来
	lastNum := fileName[len(fileName)-1:]
	return lastNum
}

// 过滤sections中的敏感字符串、
// 将含有敏感字符串的section.Text设置为空字符串
func filterDirtyText(sections []model.Section) []model.Section {

	// 过滤字符串
	if len(sections) > 0 {
		for i := len(sections) - 1; i >= 0; i-- {
			if len(sections[i].Text) > 0 {
				for _, dirtyText := range sectionDirtyTexts {
					if strings.Contains(sections[i].Text, dirtyText) {
						sections[i].Text = ""
						break
					}
				}
			}
		}
	}

	return sections
}

// 过滤sections中的section.ImageUrls
// 将sections中section.ImageUrls为空数组的section从sections中剔除
func filterEmptyImageUrls(sections []model.Section) []model.Section {

	// Filter out sections with empty image_urls
	filteredSections := make([]model.Section, 0, len(sections))
	for _, section := range sections {
		if len(section.ImageUrls) > 0 {
			filteredSections = append(filteredSections, section)
		}
	}

	return filteredSections
}
