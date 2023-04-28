/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-28 17:06:35
 * @FilePath: \image_hub\spiders\four_page.go
 * @Description: 微信公众号第4条内容抓取-表情包
 */

package spiders

import (
	"fmt"
	"image_hub/model"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	log "github.com/sirupsen/logrus"
)

type fourPage struct {
	*baseSpider
}

// NewFourPage
func NewFourPage(name string) Spider {
	return &fourPage{
		baseSpider: &baseSpider{
			Name: name,
		},
	}
}

// 解析将爬取到的数据至一个规范的结构体中
// e *colly.HTMLElement 或者  *colly.Response
func (s *fourPage) ParseData(q *queue.Queue, i interface{}, baseUrl string) (interface{}, error) {

	articleBase, err := s.baseSpider.ParseData(q, i, baseUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid type: %T, expected *colly.HTMLElement", i)
	}

	// 类型断言进行转换
	article, ok := articleBase.(*model.TblArticle)
	if !ok {
		fmt.Printf("%s failed to convert article to tblArticle", s.GetName())
		return nil, fmt.Errorf("%s failed to convert article to tblArticle", s.GetName())
	}

	e, ok := i.(*colly.HTMLElement)
	if !ok {
		return nil, fmt.Errorf("invalid type: %T, expected *colly.HTMLElement", i)
	}
	var sections []model.Section
	url := e.Request.URL.String()

	// 文章标题
	selector := "h1#activity-name"
	title := e.ChildText(selector)

	// 全部文字
	text := "🤎"

	// 所有的图片
	selector = "section section .wxw-img"
	imageUrls := e.ChildAttrs(selector, "src")

	if len(imageUrls) >= 36 {

		// 第1行文字
		// 第1组9张图
		section1ImageUrls := imageUrls[0:9]

		section1 := model.Section{
			Text:      text,
			ImageUrls: section1ImageUrls,
		}

		// 第2行文字
		// 第2组9张图
		section2ImageUrls := imageUrls[9:18]
		section2 := model.Section{
			Text:      text,
			ImageUrls: section2ImageUrls,
		}

		// 第3组文字
		// 第3组9张图
		section3ImageUrls := imageUrls[18:27]
		section3 := model.Section{
			Text:      text,
			ImageUrls: section3ImageUrls,
		}

		// 第4组文字
		// 第4组9张图
		section4ImageUrls := imageUrls[27:36]
		section4 := model.Section{
			Text:      text,
			ImageUrls: section4ImageUrls,
		}

		sections = append(sections,
			section1,
			section2,
			section3,
			section4,
		)
	} else {
		log.Warningf("imageUrls count error. : url: %s, title: %s, len(texts): %d \n", url, title, len(imageUrls))
		fmt.Printf("================ WARNING imageUrls count error. : url: %s, title: %s, len(texts): %d \n", url, title, len(imageUrls))
	}

	article.Sections = sections
	return article, nil
}
