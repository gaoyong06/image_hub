/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-25 22:21:14
 * @FilePath: \image_hub\spiders\third_page.go
 * @Description: 微信公众号第3条内容抓取-壁纸
 */

package spiders

import (
	"fmt"
	"image_hub/model"
	"image_hub/pkg/utils"
	"net/url"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

type thirdPage struct {
	Name string
}

// NewThirdPage
func NewThirdPage(name string) Spider {
	return &thirdPage{
		Name: name,
	}
}

// 获取爬虫名称
func (s *thirdPage) GetName() string {
	return s.Name
}

// 设置爬虫名称
func (s *thirdPage) SetName(name string) {
	s.Name = name
}

// 向队列追求爬取请求
func (s *thirdPage) AddReqToQueue(q *queue.Queue, i interface{}, path string) error {

	pathUrl := fmt.Sprintf("file://%s", path)

	// 解析 URL
	url, err := url.Parse(pathUrl)
	if err != nil {
		log.Errorf("url.Parse failed. err: %+v\n", err)
		return err
	}

	if _, ok := visited.Get(path); !ok {

		visited.Set(path, true)
		req := &colly.Request{
			URL:    url,
			Method: "GET",
			Ctx:    colly.NewContext(),
		}

		req.Ctx.Put(UrlTypeKey, ThirdPage)
		q.AddRequest(req)

	}
	return nil
}

// 解析将爬取到的数据至一个规范的结构体中
// e *colly.HTMLElement 或者  *colly.Response
func (s *thirdPage) ParseData(q *queue.Queue, i interface{}, baseUrl string) (interface{}, error) {

	// 解析返回html结果
	article := &model.TblArticle{}
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
	article.Title = title

	// 作者
	selector = "a#js_name"
	author := e.ChildText(selector)
	article.Author = author

	// 发布时间
	publishTime, _ := utils.GetPublishTime(e.Text)
	article.PublishTime = time.Unix(publishTime, 0)

	// <meta content="http://mp.weixin.qq.com/s?__biz=MjM5NzAyMDIwMA==&amp;mid=2653562471&amp;idx=1&amp;sn=5a209eca9a0c9d92d484dadfa516a807&amp;chksm=bd3ed1208a49583679dddb80f504983511b6bc9d63c89242dd3df68daebd587a78b8fea1afa0#rd"/>
	selector = "meta[property='og:url']"
	ogUrl := e.ChildAttr(selector, "content")
	queryParams, err := utils.GetArticleUrlQueryParams(ogUrl)
	if err != nil {
		log.Errorf("utils.GetArticleUrlQueryParams failed. ogUrl: %s,  err: %+v\n", ogUrl, err)
		return nil, err
	}
	idx := queryParams.Get("idx")
	sn := queryParams.Get("sn")
	biz := queryParams.Get("__biz")
	mid := queryParams.Get("mid")

	article.Idx = cast.ToInt(idx)
	article.Sn = sn
	article.Biz = biz
	article.Mid = cast.ToInt(mid)

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

	article.Sections = sections
	return article, nil
}

// 业务处理
// 1. 向队列追加请求
// 2. 解析数据至结构体
// 3. 保存数据 或 更新数据 或 继续下一层级的请求
// e *colly.HTMLElement 或者  *colly.Response
func (s *thirdPage) Process(q *queue.Queue, i interface{}, baseUrl string) error {

	e, ok := i.(*colly.HTMLElement)
	if !ok {
		return fmt.Errorf("%s invalid type: %T, expected *colly.HTMLElement", s.GetName(), i)
	}

	// 解析返回json结果
	article, err := s.ParseData(q, e, baseUrl)
	if err != nil {
		log.Errorf("%s ParseData failed. err: %s, url: %+v\n", s.GetName(), err, e.Request.URL.String())
		return err
	}

	// 类型断言进行转换
	tblArticle, ok := article.(*model.TblArticle)
	if ok {

		// 保存数据
		// 保存到本地article
		sn, err := tblArticle.CreateOrUpdate()
		if err != nil {

			log.Errorf("%s article.CreateOrUpdate failed. err: %s\n", s.GetName(), err)
			fmt.Printf("%s article.CreateOrUpdate failed. err: %s\n", s.GetName(), err)
			return err
		}

		fmt.Printf("%s article.CreateOrUpdate success. sn: %s\n", s.GetName(), sn)
		log.Infof("%s article.CreateOrUpdate success. sn: %s\n", s.GetName(), sn)

		// 按照多个section保存至content_service
		// TODO:调用content_service API完成批量写入

		return nil

	} else {

		fmt.Printf("%s failed to convert article to tblArticle", s.GetName())
		return fmt.Errorf("%s failed to convert article to tblArticle", s.GetName())
	}
}
