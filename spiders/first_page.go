/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-24 10:54:31
 * @FilePath: \image_hub\spiders\first_page.go
 * @Description: 微信公众号第1条内容抓取
 */

package spiders

import (
	"fmt"
	"image_hub/model"
	"image_hub/pkg/utils"
	"net/url"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
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
	url, err := url.Parse(pathUrl)
	if err != nil {
		log.Errorf("firstPage url.Parse failed. err: %+v\n", err)
		return err
	}

	if _, ok := visited.Get(pathUrl); !ok {

		visited.Set(pathUrl, true)
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
	var err error

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
	selector = "em#publish_time"
	publishTime := e.ChildText(selector)

	// 第1行文字
	section1Text := ""

	// 第1组4张图
	selector = "section:nth-child(6) p .wxw-img , p+ section > section > p .wxw-img , section:nth-child(3) section section .wxw-img"
	section1Urls := e.ChildAttrs(selector, "src")

	section1 := model.Section{
		Text:      section1Text,
		ImageUrls: section1Urls,
	}
	sections = append(sections, section1)

	// 第2行文字
	selector = "div#js_content p:nth-child(8) > span"
	section2Text := utils.FilterHTMLTags(e.ChildText(selector))

	// 第2组4张图
	selector = "section:nth-child(12) p .wxw-img , section:nth-child(10) p .wxw-img"
	section2Urls := e.ChildAttrs(selector, "src")

	section2 := model.Section{
		Text:      section2Text,
		ImageUrls: section2Urls,
	}
	sections = append(sections, section2)

	// 第3组文字
	selector = "div#js_content p:nth-child(14) > span:nth-child(3)"
	section3Text := utils.FilterHTMLTags(e.ChildText(selector))
	// 第3组4张图
	selector = "section:nth-child(18) .wxw-img , section:nth-child(16) .wxw-img"
	section3Urls := e.ChildAttrs(selector, "src")

	section3 := model.Section{
		Text:      section3Text,
		ImageUrls: section3Urls,
	}
	sections = append(sections, section3)

	// 第4组文字
	selector = "div#js_content p:nth-child(20) > span:nth-child(3)"
	section4Text := utils.FilterHTMLTags(e.ChildText(selector))
	// 第4组4张图
	selector = "section:nth-child(24) .wxw-img , p+ section section section p .wxw-img"
	section4Urls := e.ChildAttrs(selector, "src")

	section4 := model.Section{
		Text:      section4Text,
		ImageUrls: section4Urls,
	}
	sections = append(sections, section4)

	// 第5组文字
	selector = "div#js_content section:nth-child(25) > section > section > p:nth-child(2) > span:nth-child(6)"
	section5Text := utils.FilterHTMLTags(e.ChildText(selector))

	// 第5组4张图
	selector = "section:nth-child(28) .wxw-img , section:nth-child(26) .wxw-img"
	section5Urls := e.ChildAttrs(selector, "src")

	section5 := model.Section{
		Text:      section5Text,
		ImageUrls: section5Urls,
	}
	sections = append(sections, section5)

	// 第6组文字
	selector = "div#js_content section:nth-child(29) > section > section > section:nth-child(4) > span"
	section6Text := utils.FilterHTMLTags(e.ChildText(selector))

	// 第6组4张图
	selector = "section:nth-child(32) .wxw-img , section:nth-child(30) .wxw-img"
	section6Urls := e.ChildAttrs(selector, "src")

	section6 := model.Section{
		Text:      section6Text,
		ImageUrls: section6Urls,
	}
	sections = append(sections, section6)

	// 第7组文字
	selector = "div#js_content p:nth-child(34) > span"
	section7Text := utils.FilterHTMLTags(e.ChildText(selector))
	// 第7组4张图
	selector = "section:nth-child(36) .wxw-img , section:nth-child(38) .wxw-img"
	section7Urls := e.ChildAttrs(selector, "src")

	section7 := model.Section{
		Text:      section7Text,
		ImageUrls: section7Urls,
	}
	sections = append(sections, section7)

	// 第8组文字
	selector = "#js_content > section > section > section > section:nth-child(39) > p:nth-child(2) > span:nth-child(11)"
	section8Text := utils.FilterHTMLTags(e.ChildText(selector))
	// 第8组4张图
	selector = "section:nth-child(40) p .wxw-img"
	section8Urls := e.ChildAttrs(selector, "src")

	section8 := model.Section{
		Text:      section8Text,
		ImageUrls: section8Urls,
	}
	sections = append(sections, section8)

	// 第9组文字
	selector = "#js_content > section > section > section > section:nth-child(41) > p:nth-child(2) > span:nth-child(3)"
	section9Text := utils.FilterHTMLTags(e.ChildText(selector))

	// 第9组4张图
	selector = "section:nth-child(42) p .wxw-img"
	section9Urls := e.ChildAttrs(selector, "src")

	section9 := model.Section{
		Text:      section9Text,
		ImageUrls: section9Urls,
	}
	sections = append(sections, section9)

	// 第10组文字-里面有html标签
	selector = "#js_content > section > section > section > section:nth-child(43) > section:nth-child(2) > span:nth-child(6)"
	section10Text := utils.FilterHTMLTags(e.ChildText(selector))
	// 第10组4张图
	selector = "section:nth-child(44) .wxw-img , section:nth-child(45) .wxw-img , section:nth-child(46) .wxw-img , section:nth-child(47) .wxw-img"
	section10Urls := e.ChildAttrs(selector, "src")

	section10 := model.Section{
		Text:      section10Text,
		ImageUrls: section10Urls,
	}
	sections = append(sections, section10)

	// 第11组文字
	selector = "#js_content > section > section > section > section:nth-child(48) > section > section > section:nth-child(3) > span > strong > em > span > strong > em > span"
	section11Text := utils.FilterHTMLTags(e.ChildText(selector))
	// 第11组16张图
	selector = "section:nth-child(4) section .wxw-img , section:nth-child(5) section .wxw-img , section:nth-child(6) section .wxw-img , section:nth-child(7) section .wxw-img"
	section11Urls := e.ChildAttrs(selector, "src")

	section11 := model.Section{
		Text:      section11Text,
		ImageUrls: section11Urls,
	}
	sections = append(sections, section11)

	// 第12组文字(不是"你们要的",是最底部文案)
	selector = "#js_content > section > section > section > section:nth-child(48) > section > section > section:nth-child(15) > span:nth-child(2)"
	section12Text := utils.FilterHTMLTags(e.ChildText(selector))
	// 第12组16张图
	selector = "section:nth-child(10) section .wxw-img , section:nth-child(11) section .wxw-img , section:nth-child(12) section .wxw-img , section:nth-child(13) section .wxw-img"
	section12Urls := e.ChildAttrs(selector, "src")

	section12 := model.Section{
		Text:      section12Text,
		ImageUrls: section12Urls,
	}
	sections = append(sections, section12)

	article.Title = title
	article.Author = author
	article.PublishTime, err = utils.StringToTime(publishTime)
	if err != nil {
		log.Errorf("firstPage ParseData utils.StringToTime(publishTime) failed. err: %s, publishTime: %s\n", err, publishTime)
	}

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

	fmt.Printf("firstPage Process complete. article: %#v", article)

	// // 保存数据
	// modelDetailId, err := tblModel.CreateOrUpdate()
	// if err != nil {
	// 	log.Errorf("CarParamSpider TblCarParam Create failed. err: %s\n", err)
	// 	return err
	// }
	// log.Infof("CarParam create success. modelDetailId: %d\n", modelDetailId)
	return nil
}
