/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-28 17:30:00
 * @FilePath: \image_hub\spiders\first_page.go
 * @Description: 微信公众号第1条内容抓取-头像
 */

package spiders

import (
	"fmt"
	"image_hub/model"
	"image_hub/pkg/utils"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/samber/lo"

	log "github.com/sirupsen/logrus"
)

type firstPage struct {
	*baseSpider
}

// NewFirstPage
func NewFirstPage(name string) Spider {
	return &firstPage{
		baseSpider: &baseSpider{
			Name: name,
		},
	}
}

// 解析将爬取到的数据至一个规范的结构体中
// e *colly.HTMLElement 或者  *colly.Response
func (s *firstPage) ParseData(q *queue.Queue, i interface{}, baseUrl string) (interface{}, error) {

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

	// 所有的文字
	// 下去取文字的地方有个bug,  本来是"🔥 𝑳𝒐𝒗𝒆 𝒎𝒆 𝒆𝒗𝒆𝒓𝒚𝒅𝒂𝒚",最后取到的是 "❤️\u200d🔥 𝑳𝒐𝒗𝒆 𝒎𝒆 𝒆𝒗𝒆𝒓𝒚𝒅𝒂𝒚"
	// 文档地址：file:///D:/work/wechat_download_data/html/Dump-0421-11-15-39/20220526_111900_1.html
	selector = "section, p"

	var texts []string
	e.ForEach(selector, func(i int, h *colly.HTMLElement) {

		// fmt.Printf("============ url: %s, title: %s, h.Text: h.Text %+v\n", url, title, h.Text)
		texts = append(texts, h.Text)
	})

	// fmt.Printf("================ 原始字符串数组: url: %s, title: %s, len(texts): %d,  texts: %#v\n", url, title, len(texts), texts)

	// 遍历texts，从后向前遍历，如果前面的项的字符串中，完整包含了后面项的字符串，则将前面项的字符串，设置为空字符串
	for i := len(texts) - 1; i >= 0; i-- {

		if len(texts[i]) > 0 {
			for j := 0; j < i; j++ {
				if len(texts[i]) > 0 {
					if strings.Contains(texts[j], texts[i]) {
						texts[j] = ""
					}
				}
			}
		}
	}
	// fmt.Printf("================ 字符串去重后: url: %s, title: %s, len(texts): %d,  texts: %#v\n", url, title, len(texts), texts)

	// 过滤字符串
	for i := len(texts) - 1; i >= 0; i-- {
		if len(texts[i]) > 0 {
			for _, dirtyText := range sectionDirtyTexts {
				if strings.Contains(texts[i], dirtyText) {
					texts[i] = ""
					break
				}
			}
		}
	}
	// fmt.Printf("================ 字符串过滤后: url: %s, title: %s, len(texts): %d,  texts: %#v\n", url, title, len(texts), texts)

	// 将前后连续的字符串使用\n连接为一个,被连接的设置为空字符串
	texts = utils.JoinAdjacentStrings(texts)
	// fmt.Printf("================ 字符串连接后: url: %s, title: %s, len(texts): %d,  texts: %#v\n", url, title, len(texts), texts)

	// 过滤掉所有的空字符串
	texts = lo.Filter(texts, func(text string, idx int) bool {

		text = strings.ReplaceAll(text, "\n", "")
		text = strings.ReplaceAll(text, " ", "")
		if len(text) == 0 {
			return false
		} else {
			return true
		}
	})

	// fmt.Printf("================ 过滤掉所有的空字符串后: url: %s, title: %s, len(texts): %d,  texts: %#v\n", url, title, len(texts), texts)

	if len(texts) != 11 {

		log.Warningf("texts count error. : url: %s, title: %s, len(texts): %d,  texts: %#v\n", url, title, len(texts), texts)
		fmt.Printf("================ WARNING texts count error. : url: %s, title: %s, len(texts): %d,  texts: %#v\n", url, title, len(texts), texts)
	}

	// 不足11个，补全为11个
	for len(texts) < 11 {
		texts = append(texts, texts[len(texts)-1])
	}

	// fmt.Printf("================ 最终使用的texts: url: %s, title: %s, len(texts): %d,  texts: %#v\n", url, title, len(texts), texts)

	// 所有的图片
	selector = ".wxw-img"
	imageUrls := e.ChildAttrs(selector, "src")

	// 一共有72张图
	if len(imageUrls) >= 72 {

		// 删掉最后一张图
		imageUrls = imageUrls[:len(imageUrls)-1]

		// 第1行文字
		section1Text := ""
		// 第1组4张图
		section1ImageUrls := imageUrls[0:4]

		section1 := model.Section{
			Text:      section1Text,
			ImageUrls: section1ImageUrls,
		}

		// 第2行文字
		section2Text := texts[0]
		// 第2组4张图
		section2ImageUrls := imageUrls[4:8]
		section2 := model.Section{
			Text:      section2Text,
			ImageUrls: section2ImageUrls,
		}

		// 第3组文字
		section3Text := texts[1]
		// 第3组4张图
		section3ImageUrls := imageUrls[8:12]
		section3 := model.Section{
			Text:      section3Text,
			ImageUrls: section3ImageUrls,
		}

		// 第4组文字
		section4Text := texts[2]
		// 第4组4张图
		section4ImageUrls := imageUrls[12:16]
		section4 := model.Section{
			Text:      section4Text,
			ImageUrls: section4ImageUrls,
		}

		// 第5组文字
		section5Text := texts[3]
		// 第5组4张图
		section5ImageUrls := imageUrls[16:20]
		section5 := model.Section{
			Text:      section5Text,
			ImageUrls: section5ImageUrls,
		}

		// 第6组文字
		section6Text := texts[4]
		// 第6组4张图
		section6ImageUrls := imageUrls[20:24]
		section6 := model.Section{
			Text:      section6Text,
			ImageUrls: section6ImageUrls,
		}

		// 第7组文字
		section7Text := texts[5]
		// 第7组4张图
		section7ImageUrls := imageUrls[24:28]
		section7 := model.Section{
			Text:      section7Text,
			ImageUrls: section7ImageUrls,
		}

		// 第8组文字
		section8Text := texts[6]
		// 第8组4张图
		section8ImageUrls := imageUrls[28:32]
		section8 := model.Section{
			Text:      section8Text,
			ImageUrls: section8ImageUrls,
		}

		// 第9组文字
		section9Text := texts[7]
		// 第9组4张图
		section9ImageUrls := imageUrls[32:36]
		section9 := model.Section{
			Text:      section9Text,
			ImageUrls: section9ImageUrls,
		}

		// 第10组文字
		section10Text := texts[8]
		// 第10组4张图
		section10ImageUrls := imageUrls[36:40]
		section10 := model.Section{
			Text:      section10Text,
			ImageUrls: section10ImageUrls,
		}

		// 第11组文字
		section11Text := texts[9]
		// 第11组16张图
		section11ImageUrls := imageUrls[40:56]
		section11 := model.Section{
			Text:      section11Text,
			ImageUrls: section11ImageUrls,
		}

		// 第12组文字
		section12Text := texts[10]
		// 第12组16张图
		section12ImageUrls := imageUrls[56:72]
		section12 := model.Section{
			Text:      section12Text,
			ImageUrls: section12ImageUrls,
		}

		sections = append(sections,
			section1,
			section2,
			section3,
			section4,
			section5,
			section6,
			section7,
			section8,
			section9,
			section10,
			section11,
			section12,
		)
	} else {
		log.Warningf("imageUrls count error. : url: %s, title: %s, len(texts): %d \n", url, title, len(imageUrls))
		fmt.Printf("================ WARNING imageUrls count error. : url: %s, title: %s, len(texts): %d \n", url, title, len(imageUrls))
	}

	article.Sections = sections
	return article, nil
}
