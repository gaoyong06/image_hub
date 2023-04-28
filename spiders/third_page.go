/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-28 17:03:26
 * @FilePath: \image_hub\spiders\third_page.go
 * @Description: 微信公众号第3条内容抓取-壁纸
 */

package spiders

import (
	"fmt"
	"image_hub/model"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	log "github.com/sirupsen/logrus"
)

type thirdPage struct {
	*baseSpider
}

// NewThirdPage
func NewThirdPage(name string) Spider {
	return &thirdPage{
		baseSpider: &baseSpider{
			Name: name,
		},
	}
}

// 解析将爬取到的数据至一个规范的结构体中
// e *colly.HTMLElement 或者  *colly.Response
func (s *thirdPage) ParseData(q *queue.Queue, i interface{}, baseUrl string) (interface{}, error) {

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
	// 文字有两种
	//  1. 最后一行图片下面一行文字
	//  2. 其他都是 🌷 🤍 🌷
	selector = "section section span"
	lastText := e.ChildText(selector)

	text := "🌷 🤍 🌷"

	// 所有的图片
	selector = "section section .wxw-img"
	imageUrls := e.ChildAttrs(selector, "src")

	if len(imageUrls) >= 42 {
		// 第1行文字
		// 第1组6张图
		section1ImageUrls := imageUrls[0:6]

		section1 := model.Section{
			Text:      text,
			ImageUrls: section1ImageUrls,
		}

		// 第2行文字
		// 第2组6张图
		section2ImageUrls := imageUrls[6:12]
		section2 := model.Section{
			Text:      text,
			ImageUrls: section2ImageUrls,
		}

		// 第3组文字
		// 第3组6张图
		section3ImageUrls := imageUrls[12:18]
		section3 := model.Section{
			Text:      text,
			ImageUrls: section3ImageUrls,
		}

		// 第4组文字
		// 第4组6张图
		section4ImageUrls := imageUrls[18:24]
		section4 := model.Section{
			Text:      text,
			ImageUrls: section4ImageUrls,
		}

		// 第5组文字
		// 第5组6张图
		section5ImageUrls := imageUrls[24:30]
		section5 := model.Section{
			Text:      text,
			ImageUrls: section5ImageUrls,
		}

		// 第6组文字
		// 第6组6张图
		section6ImageUrls := imageUrls[30:36]
		section6 := model.Section{
			Text:      text,
			ImageUrls: section6ImageUrls,
		}

		// 第7组文字
		// 第7组6张图
		section7ImageUrls := imageUrls[36:42]
		section7 := model.Section{
			Text:      lastText,
			ImageUrls: section7ImageUrls,
		}

		sections = append(sections,
			section1,
			section2,
			section3,
			section4,
			section5,
			section6,
			section7,
		)

	} else {
		log.Warningf("imageUrls count error. : url: %s, title: %s, len(texts): %d \n", url, title, len(imageUrls))
		fmt.Printf("================ WARNING imageUrls count error. : url: %s, title: %s, len(texts): %d \n", url, title, len(imageUrls))
	}

	article.Sections = sections
	return article, nil
}
