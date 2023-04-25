/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-25 18:01:20
 * @FilePath: \image_hub\spiders\first_page.go
 * @Description: 微信公众号第1条内容抓取
 */

package spiders

import (
	"fmt"
	"image_hub/model"
	"image_hub/pkg/utils"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

type firstPage struct {
	Name string
}

// NewFirstPage
func NewFirstPage(name string) Spider {
	return &firstPage{
		Name: name,
	}
}

// 获取爬虫名称
func (s *firstPage) GetName() string {
	return s.Name
}

// 设置爬虫名称
func (s *firstPage) SetName(name string) {
	s.Name = name
}

// 向队列追求爬取请求
func (s *firstPage) AddReqToQueue(q *queue.Queue, i interface{}, path string) error {

	pathUrl := fmt.Sprintf("file://%s", path)

	// 解析 URL
	url, err := url.Parse(pathUrl)
	if err != nil {
		log.Errorf("firstPage url.Parse failed. err: %+v\n", err)
		return err
	}

	if _, ok := visited.Get(path); !ok {

		visited.Set(path, true)
		req := &colly.Request{
			URL:    url,
			Method: "GET",
			Ctx:    colly.NewContext(),
		}

		req.Ctx.Put(UrlTypeKey, FirstPage)
		q.AddRequest(req)

	}
	return nil
}

// 解析将爬取到的数据至一个规范的结构体中
// e *colly.HTMLElement 或者  *colly.Response
func (s *firstPage) ParseData(q *queue.Queue, i interface{}, baseUrl string) (interface{}, error) {

	// 解析返回html结果
	article := &model.Article{}
	var selector string
	var sections []model.Section
	// var err error

	e, ok := i.(*colly.HTMLElement)
	if !ok {
		return nil, fmt.Errorf("invalid type: %T, expected *colly.HTMLElement", i)
	}

	// 文章标题
	selector = "h1#activity-name"
	title := e.ChildText(selector)

	// 作者
	selector = "a#js_name"
	author := e.ChildText(selector)

	// 发布时间
	publishTime, _ := utils.GetPublishTime(e.Text)

	// 所有的文字
	// 下去取文字的地方有个bug,  本来是"🔥 𝑳𝒐𝒗𝒆 𝒎𝒆 𝒆𝒗𝒆𝒓𝒚𝒅𝒂𝒚",最后取到的是 "❤️\u200d🔥 𝑳𝒐𝒗𝒆 𝒎𝒆 𝒆𝒗𝒆𝒓𝒚𝒅𝒂𝒚"
	// 文档地址：file:///D:/work/wechat_download_data/html/Dump-0421-11-15-39/20220526_111900_1.html
	selector = "section span:not(.audio_area,  .audio_area  *), p span"
	var textsStr string
	e.ForEach(selector, func(i int, h *colly.HTMLElement) {

		nodes := h.DOM.Nodes
		for _, n := range nodes {

			if n.Attr[0].Key == "style" {

				// 中间各个区块的名称
				if strings.Contains(n.Attr[0].Val, "text-align: center;") {

					textsStr = textsStr + h.Text
					break
				}

				// 倒数第3个文字：真人头像，倒数第2个文字：你们要的
				if strings.Contains(n.Attr[0].Val, "text-decoration: underline") {
					textsStr = textsStr + h.Text
					break
				}
			} else {
				textsStr = textsStr + "\t"
				break
			}

			// 最后1行文字：我好想你啊 这句话无论谁和我说起 我都会想要掉眼泪 我就觉得被人惦记 真好啊
			// 无法通过 style="font-size: 12px;color: rgb(73, 73, 73);font-family: Optima-Regular, PingFangTC-light;" 这匹配, 这个样式不是固定的
			// 例如: 下面这个
			// https://mp.weixin.qq.com/s/nXAfWugJouIEQ4hhbAcStg
			// 目前通过节点的索引号，和文字长度来判断
			if i > 60 && len(h.Text) > 20 {
				textsStr = textsStr + h.Text
				break
			} else {
				textsStr = textsStr + "\t"
				break
			}
		}
	})

	texts := strings.FieldsFunc(textsStr, func(r rune) bool {
		return r == '\t'
	})

	// 过滤掉"你们要的"
	uselessStr := "你们要的"
	texts = lo.Filter(texts, func(val string, idx int) bool {

		return !strings.Contains(val, uselessStr)
	})

	// 所有的图片
	// .wxw-img

	// // 第1行文字
	// section1Text := ""

	// // 第1组4张图
	// selector = "section:nth-child(6) p .wxw-img , p+ section > section > p .wxw-img , section:nth-child(3) section section .wxw-img"
	// section1Urls := e.ChildAttrs(selector, "src")
	// fmt.Printf("第1组4张图 %+v\n", section1Urls)

	// section1 := model.Section{
	// 	Text:      section1Text,
	// 	ImageUrls: section1Urls,
	// }
	// sections = append(sections, section1)

	// // 第2行文字
	// selector = "#js_content > section:nth-child(4)"
	// section2Text := e.ChildText(selector)
	// fmt.Printf("第2组文字 %+v\n", section2Text)

	// // 第2组4张图

	// // selector = "section:nth-child(12) p .wxw-img , section:nth-child(10) p .wxw-img"
	// // #js_content > section:nth-child(6)
	// // #js_content > section:nth-child(6) .wxw-img, #js_content > section:nth-child(8) .wxw-img
	// // selector = "#js_content > section:nth-child(6) > section:nth-child(1) > section > section > section > img"
	// selector = "#js_content > section:nth-child(6) .wxw-img, #js_content > section:nth-child(8) .wxw-img"
	// section2Urls := e.ChildAttrs(selector, "src")
	// fmt.Printf("第2组4张图 %+v\n", section2Urls)

	// return nil, nil

	// section2 := model.Section{
	// 	Text:      section2Text,
	// 	ImageUrls: section2Urls,
	// }
	// sections = append(sections, section2)

	// // 第3组文字
	// selector = "div#js_content p:nth-child(14) > span:nth-child(3)"
	// section3Text := utils.FilterHTMLTags(e.ChildText(selector))
	// fmt.Printf("第3组文字 %+v\n", section3Text)
	// // 第3组4张图
	// selector = "section:nth-child(18) .wxw-img , section:nth-child(16) .wxw-img"
	// section3Urls := e.ChildAttrs(selector, "src")
	// fmt.Printf("第3组4张图 %+v\n", section3Urls)

	// section3 := model.Section{
	// 	Text:      section3Text,
	// 	ImageUrls: section3Urls,
	// }
	// sections = append(sections, section3)

	// // 第4组文字
	// selector = "div#js_content p:nth-child(20) > span:nth-child(3)"
	// section4Text := utils.FilterHTMLTags(e.ChildText(selector))
	// fmt.Printf("第4组文字 %+v\n", section4Text)
	// // 第4组4张图
	// selector = "section:nth-child(24) .wxw-img , p+ section section section p .wxw-img"
	// section4Urls := e.ChildAttrs(selector, "src")
	// fmt.Printf("第4组4张图 %+v\n", section4Urls)

	// section4 := model.Section{
	// 	Text:      section4Text,
	// 	ImageUrls: section4Urls,
	// }
	// sections = append(sections, section4)

	// // 第5组文字
	// selector = "div#js_content section:nth-child(25) > section > section > p:nth-child(2) > span:nth-child(6)"
	// section5Text := utils.FilterHTMLTags(e.ChildText(selector))
	// fmt.Printf("第5组文字 %+v\n", section5Text)

	// // 第5组4张图
	// selector = "section:nth-child(28) .wxw-img , section:nth-child(26) .wxw-img"
	// section5Urls := e.ChildAttrs(selector, "src")
	// fmt.Printf("第5组4张图 %+v\n", section5Urls)

	// section5 := model.Section{
	// 	Text:      section5Text,
	// 	ImageUrls: section5Urls,
	// }
	// sections = append(sections, section5)

	// // 第6组文字
	// selector = "div#js_content section:nth-child(29) > section > section > section:nth-child(4) > span"
	// section6Text := utils.FilterHTMLTags(e.ChildText(selector))
	// fmt.Printf("第6组文字 %+v\n", section6Text)

	// // 第6组4张图
	// selector = "section:nth-child(32) .wxw-img , section:nth-child(30) .wxw-img"
	// section6Urls := e.ChildAttrs(selector, "src")
	// fmt.Printf("第6组4张图 %+v\n", section6Urls)

	// section6 := model.Section{
	// 	Text:      section6Text,
	// 	ImageUrls: section6Urls,
	// }
	// sections = append(sections, section6)

	// // 第7组文字
	// selector = "div#js_content p:nth-child(34) > span"
	// section7Text := utils.FilterHTMLTags(e.ChildText(selector))
	// fmt.Printf("第7组文字 %+v\n", section7Text)
	// // 第7组4张图
	// selector = "section:nth-child(36) .wxw-img , section:nth-child(38) .wxw-img"
	// section7Urls := e.ChildAttrs(selector, "src")
	// fmt.Printf("第7组4张图 %+v\n", section7Urls)

	// section7 := model.Section{
	// 	Text:      section7Text,
	// 	ImageUrls: section7Urls,
	// }
	// sections = append(sections, section7)

	// // 第8组文字
	// selector = "#js_content > section > section > section > section:nth-child(39) > p:nth-child(2) > span:nth-child(11)"
	// section8Text := utils.FilterHTMLTags(e.ChildText(selector))
	// fmt.Printf("第8组文字 %+v\n", section8Text)
	// // 第8组4张图
	// selector = "section:nth-child(40) p .wxw-img"
	// section8Urls := e.ChildAttrs(selector, "src")
	// fmt.Printf("第8组4张图 %+v\n", section8Urls)

	// section8 := model.Section{
	// 	Text:      section8Text,
	// 	ImageUrls: section8Urls,
	// }
	// sections = append(sections, section8)

	// // 第9组文字
	// selector = "#js_content > section > section > section > section:nth-child(41) > p:nth-child(2) > span:nth-child(3)"
	// section9Text := utils.FilterHTMLTags(e.ChildText(selector))
	// fmt.Printf("第9组文字 %+v\n", section9Text)

	// // 第9组4张图
	// selector = "section:nth-child(42) p .wxw-img"
	// section9Urls := e.ChildAttrs(selector, "src")
	// fmt.Printf("第9组4张图 %+v\n", section9Urls)

	// section9 := model.Section{
	// 	Text:      section9Text,
	// 	ImageUrls: section9Urls,
	// }
	// sections = append(sections, section9)

	// // 第10组文字-里面有html标签
	// selector = "#js_content > section > section > section > section:nth-child(43) > section:nth-child(2) > span:nth-child(6)"
	// section10Text := utils.FilterHTMLTags(e.ChildText(selector))
	// fmt.Printf("第10组文字 %+v\n", section10Text)
	// // 第10组4张图
	// selector = "section:nth-child(44) .wxw-img , section:nth-child(45) .wxw-img , section:nth-child(46) .wxw-img , section:nth-child(47) .wxw-img"
	// section10Urls := e.ChildAttrs(selector, "src")
	// fmt.Printf("第10组4张图 %+v\n", section10Urls)

	// section10 := model.Section{
	// 	Text:      section10Text,
	// 	ImageUrls: section10Urls,
	// }
	// sections = append(sections, section10)

	// // 第11组文字
	// selector = "#js_content > section > section > section > section:nth-child(48) > section > section > section:nth-child(3) > span > strong > em > span > strong > em > span"
	// section11Text := utils.FilterHTMLTags(e.ChildText(selector))
	// fmt.Printf("第11组文字 %+v\n", section11Text)
	// // 第11组16张图
	// selector = "section:nth-child(4) section .wxw-img , section:nth-child(5) section .wxw-img , section:nth-child(6) section .wxw-img , section:nth-child(7) section .wxw-img"
	// section11Urls := e.ChildAttrs(selector, "src")
	// fmt.Printf("第11组16张图 %+v\n", section11Urls)

	// section11 := model.Section{
	// 	Text:      section11Text,
	// 	ImageUrls: section11Urls,
	// }
	// sections = append(sections, section11)

	// // 第12组文字(不是"你们要的",是最底部文案)
	// selector = "#js_content > section > section > section > section:nth-child(48) > section > section > section:nth-child(15) > span:nth-child(2)"
	// section12Text := utils.FilterHTMLTags(e.ChildText(selector))
	// fmt.Printf("第12组文字 %+v\n", section12Text)
	// // 第12组16张图
	// selector = "section:nth-child(10) section .wxw-img , section:nth-child(11) section .wxw-img , section:nth-child(12) section .wxw-img , section:nth-child(13) section .wxw-img"
	// section12Urls := e.ChildAttrs(selector, "src")
	// fmt.Printf("第12组16张图 %+v\n", section12Urls)

	// section12 := model.Section{
	// 	Text:      section12Text,
	// 	ImageUrls: section12Urls,
	// }
	// sections = append(sections, section12)

	article.Title = title
	article.Author = author
	article.PublishTime = time.Unix(publishTime, 0)

	article.Sections = sections
	return article, nil
}

// 业务处理
// 1. 向队列追加请求
// 2. 解析数据至结构体
// 3. 保存数据 或 更新数据 或 继续下一层级的请求
// e *colly.HTMLElement 或者  *colly.Response
func (s *firstPage) Process(q *queue.Queue, i interface{}, baseUrl string) error {

	e, ok := i.(*colly.HTMLElement)
	if !ok {
		return fmt.Errorf("invalid type: %T, expected *colly.HTMLElement", i)
	}

	// 解析返回json结果
	article, err := s.ParseData(q, e, baseUrl)
	if err != nil {
		log.Errorf("firstPage ParseData failed. err: %s, url: %+v\n", err, e.Request.URL.String())
		return err
	}

	log.Infof("firstPage Process complete. article: %#v", article)
	// fmt.Printf("firstPage Process complete. article: %#v", article)

	// // 保存数据
	// modelDetailId, err := tblModel.CreateOrUpdate()
	// if err != nil {
	// 	log.Errorf("CarParamSpider TblCarParam Create failed. err: %s\n", err)
	// 	return err
	// }
	// log.Infof("CarParam create success. modelDetailId: %d\n", modelDetailId)
	return nil
}
