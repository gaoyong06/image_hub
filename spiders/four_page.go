/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-24 10:57:37
 * @FilePath: \image_hub\spiders\four_page.go
 * @Description: 微信公众号第4条内容抓取
 */

package spiders

import (
	"github.com/gocolly/colly/v2/queue"
)

type fourPage struct {
	Name string
}

// NewFourPage
func NewFourPage(name string) Spider {
	return &fourPage{
		Name: name,
	}
}

// 获取爬虫名称
func (s *fourPage) GetName() string {
	return s.Name
}

// 设置爬虫名称
func (s *fourPage) SetName(name string) {
	s.Name = name
}

// 向队列追求爬取请求
func (s *fourPage) AddReqToQueue(q *queue.Queue, i interface{}, baseUrl string) error {

	return nil
}

// 解析将爬取到的数据至一个规范的结构体中
// e *colly.HTMLElement 或者  *colly.Response
func (s *fourPage) ParseData(q *queue.Queue, i interface{}, baseUrl string) (interface{}, error) {

	return nil, nil
}

// 业务处理
// 1. 向队列追加请求
// 2. 解析数据至结构体
// 3. 保存数据 或 更新数据 或 继续下一层级的请求
// e *colly.HTMLElement 或者  *colly.Response
func (s *fourPage) Process(q *queue.Queue, i interface{}, baseUrl string) error {

	return nil
}
