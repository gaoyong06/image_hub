/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-21 19:20:23
 * @FilePath: \image_hub\spiders\first_page.go
 * @Description: 微信公众号第1条内容抓取
 */

package spiders

import (
	"fmt"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"google.golang.org/appengine/log"
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
func (s *firstPage) AddReqToQueue(q *queue.Queue, e *colly.HTMLElement, path string) error {

	url := fmt.Sprintf("file://%s", path)

	if _, ok := visited.Get(url); !ok {

		visited.Set(url, true)
		req, err := configPageHtmlElement.Request.New("GET", url, nil)
		if err != nil {
			log.Errorf("UsedCarOptionSpider AddReqToQueue Request.New failed. err: %s\n", err)
			return err
		}

		req.Ctx.Put(UrlTypeKey, UrlTypeCarParamApi)
		q.AddRequest(req)
	}
	return nil
}

// 解析将爬取到的数据至一个规范的结构体中
// e *colly.HTMLElement 或者  *colly.Response
func (s *firstPage) ParseData(q *queue.Queue, r interface{}, baseUrl string) (interface{}, error) {

	// 解析返回json结果
	// res := &model.CarParamRes{}
	// err := json.Unmarshal(r.Body, res)
	// if err != nil {

	// 	fmt.Println("err = ", err)
	// 	return nil, err
	// }

	// if res.ReturnCode == 0 {

	// 	paramTypeItems := res.Result.ParamTypeItems
	// 	if len(paramTypeItems) > 0 {

	// 		tblModel := model.GetTblCarParam()
	// 		tblModel.ModelDetailId = res.Result.SpecId
	// 		tblModel.CarParamRes = res
	// 		tblModel.ResBody = string(r.Body)
	// 		tblModel.OriginUrl = r.Request.URL.String()
	// 		return tblModel, nil
	// 	}

	// 	return nil, fmt.Errorf("CarParamSpider TblCarParam Create failed. paramTypeItems is empty")
	// }

	// err = fmt.Errorf("CarParamSpider ParseData failed. returncode:%d, message: %+v, url: %+v", res.ReturnCode, res.Message, r.Request.URL.String())
	return nil, nil
}

// 业务处理
// 1. 向队列追加请求
// 2. 解析数据至结构体
// 3. 保存数据 或 更新数据 或 继续下一层级的请求
// e *colly.HTMLElement 或者  *colly.Response
func (s *firstPage) Process(q *queue.Queue, r interface{}, baseUrl string) error {

	// 解析返回json结果
	// tblModel, err := s.ParseData(q, r, baseUrl)
	// if err != nil {
	// 	log.Errorf("CarParamSpider ParseData failed. err: %s, url: %+v\n", err, r.Request.URL.String())
	// 	return err
	// }

	// // 保存数据
	// modelDetailId, err := tblModel.CreateOrUpdate()
	// if err != nil {
	// 	log.Errorf("CarParamSpider TblCarParam Create failed. err: %s\n", err)
	// 	return err
	// }
	// log.Infof("CarParam create success. modelDetailId: %d\n", modelDetailId)
	return nil
}
