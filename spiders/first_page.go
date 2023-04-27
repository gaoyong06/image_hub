/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-27 11:07:03
 * @FilePath: \image_hub\spiders\first_page.go
 * @Description: 微信公众号第1条内容抓取-头像
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
	"github.com/spf13/cast"
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

		req.Ctx.Put(UrlTypeKey, FirstPage)
		q.AddRequest(req)

	}
	return nil
}

// 解析将爬取到的数据至一个规范的结构体中
// e *colly.HTMLElement 或者  *colly.Response
func (s *firstPage) ParseData(q *queue.Queue, i interface{}, baseUrl string) (interface{}, error) {

	// 解析返回html结果
	article := &model.TblArticle{}
	var selector string
	var sections []model.Section

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

	article.LocalPath = e.Request.URL.String()

	// 所有的文字
	// 下去取文字的地方有个bug,  本来是"🔥 𝑳𝒐𝒗𝒆 𝒎𝒆 𝒆𝒗𝒆𝒓𝒚𝒅𝒂𝒚",最后取到的是 "❤️\u200d🔥 𝑳𝒐𝒗𝒆 𝒎𝒆 𝒆𝒗𝒆𝒓𝒚𝒅𝒂𝒚"
	// 文档地址：file:///D:/work/wechat_download_data/html/Dump-0421-11-15-39/20220526_111900_1.html
	// selector = ".wxw-img~ span"
	// selector = "section, p, span"

	selector = "p[style*='text-align: center']>span, section[style*='text-align: center'] > section > span, section[style*='text-align: center'] > span, span[style*='text-align: center'], span[style*='text-decoration: underline;']"
	var textsStr string

	e.ForEach(selector, func(i int, h *colly.HTMLElement) {

		pText := h.Text
		if len(h.Text) > 0 {
			textsStr = textsStr + pText + "\n"
		} else {
			textsStr = textsStr + "\t"
		}

	})

	textsStr = strings.TrimRight(textsStr, "\n")

	texts := strings.FieldsFunc(textsStr, func(r rune) bool {
		return r == '\t'
	})

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

	// 过滤掉"你们要的"
	uselessStr := "\n你们要的"
	lo.ForEach(texts, func(text string, idx int) {

		texts[idx] = strings.Trim(text, "\n")
		texts[idx] = strings.Replace(texts[idx], uselessStr, "", -1)
	})

	// 所有的图片
	selector = ".wxw-img"
	imageUrls := e.ChildAttrs(selector, "src")

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
