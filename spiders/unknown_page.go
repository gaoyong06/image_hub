/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-28 16:30:26
 * @FilePath: \image_hub\spiders\unknown_page.go
 * @Description: 微信公众号第1条内容抓取-不处理的页面
 */

package spiders

import (
	"github.com/gocolly/colly/v2/queue"
)

type unknownPage struct {
	Name string
}

// NewUnknownPage
func NewUnknownPage(name string) Spider {
	return &unknownPage{
		Name: name,
	}
}

// 获取爬虫名称
func (s *unknownPage) GetName() string {
	return s.Name
}

// 设置爬虫名称
func (s *unknownPage) SetName(name string) {
	s.Name = name
}

// 向队列追求爬取请求
func (s *unknownPage) AddReqToQueue(q *queue.Queue, i interface{}, path string) error {

	return nil
}

// 解析将爬取到的数据至一个规范的结构体中
// e *colly.HTMLElement 或者  *colly.Response
func (s *unknownPage) ParseData(q *queue.Queue, i interface{}, baseUrl string) (interface{}, error) {

	return nil, nil
}

// 业务处理
// 1. 向队列追加请求
// 2. 解析数据至结构体
// 3. 保存数据 或 更新数据 或 继续下一层级的请求
// e *colly.HTMLElement 或者  *colly.Response
func (s *unknownPage) Process(a Spider, q *queue.Queue, i interface{}, baseUrl string) error {

	return nil
}
