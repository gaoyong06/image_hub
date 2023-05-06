/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-05-05 10:01:07
 * @FilePath: \image_hub\spiders\first_page.go
 * @Description: 微信公众号第1条内容抓取-头像
 */

package spiders

import (
	"encoding/json"
	"fmt"
	"image_hub/model"
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"golang.org/x/net/html"
)

type onePage struct {
	*baseSpider
}

// NewOnePage
func NewOnePage(name string) Spider {
	return &onePage{
		baseSpider: &baseSpider{
			Name: name,
		},
	}
}

// 解析将爬取到的数据至一个规范的结构体中
// e *colly.HTMLElement 或者  *colly.Response
func (s *onePage) ParseData(q *queue.Queue, i interface{}, baseUrl string) (interface{}, error) {

	articleBase, err := s.baseSpider.ParseData(q, i, baseUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid type: %T, expected *colly.HTMLElement", i)
	}

	// Type assertion to convert i to *colly.HTMLElement
	e, ok := i.(*colly.HTMLElement)
	if !ok {
		return nil, fmt.Errorf("invalid type: %T, expected *colly.HTMLElement", i)
	}

	// file://D:/work/wechat_download_data/html/test4/20220810_111900_1.html
	url := e.Request.URL.String()

	// 文章标题
	selector := "h1#activity-name"
	title := e.ChildText(selector)

	// 微信号
	selector = ".profile_meta_value"
	wechatId := e.ChildTexts(selector)[0]

	fmt.Printf("================== wechatId: %s,  title: %s, url: %s ================", wechatId, title, url)

	// Get the HTML byte slice of the e element
	htmlBytes := e.Response.Body
	htmlStr := string(htmlBytes)

	// Parse the HTML string to extract the sections
	sections := ParseSectionsFromHTML(htmlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sections from HTML: %v", err)
	}

	// 调用每个微信号及其内容索引的自定义方法
	fileIdx := getFileName(url)
	funcKey := fmt.Sprintf("%s%s", wechatId, fileIdx)
	sections = runFunc(funcKey, sections)

	// 将sections以json格式打印出来
	sectionsJson, err := json.Marshal(sections)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sections to json: %v", err)
	}
	fmt.Printf("\n\n=================== Sections JSON======================\n\n%s\n", sectionsJson)

	// Update the article object with the extracted sections
	article, ok := articleBase.(*model.TblArticle)
	if !ok {
		fmt.Printf("%s failed to convert article to tblArticle", s.GetName())
		return nil, fmt.Errorf("%s failed to convert article to tblArticle", s.GetName())
	}

	article.Sections = sections
	return article, nil

}

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
