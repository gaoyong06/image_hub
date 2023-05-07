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

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
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
	var wechatId string
	selector = ".profile_meta_value"
	profileMetaValues := e.ChildTexts(selector)

	if len(profileMetaValues) == 0 {
		panic(fmt.Sprintf("html class .profile_meta_value element is empty. title: %s, url: %s", title, url))
	} else {
		wechatId = profileMetaValues[0]
	}

	fmt.Printf("================== wechatId: %s,  title: %s, url: %s ================", wechatId, title, url)

	// Get the HTML byte slice of the e element
	htmlBytes := e.Response.Body
	htmlStr := string(htmlBytes)

	// Call the filterImages function and update the htmlString with the new image tags
	// filterImages, err := FilterImagesFromHTML(htmlStr)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to filter images: %v", err)
	// }

	// for _, badImg := range badImgs {
	// 	htmlStr = strings.Replace(htmlStr, fmt.Sprintf(`src="%s"`, badImg), fmt.Sprintf(`src="%s" class="bad-image"`, badImg), -1)
	// }

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

	// fmt.Printf("\n\n=================== filterImages======================\n\n%s\n", filterImages)

	panic("======================================== STOP ================================================")

	// Update the article object with the extracted sections
	article, ok := articleBase.(*model.TblArticle)
	if !ok {
		fmt.Printf("%s failed to convert article to tblArticle", s.GetName())
		return nil, fmt.Errorf("%s failed to convert article to tblArticle", s.GetName())
	}

	article.Sections = sections
	return article, nil

}
