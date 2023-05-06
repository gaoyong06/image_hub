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
	*baseSpider
}

// NewUnknownPage
func NewUnknownPage(name string) Spider {
	return &unknownPage{
		baseSpider: &baseSpider{
			Name: name,
		},
	}
}

// 向队列追求爬取请求
// 不规范的文案，暂时不做处理
func (s *unknownPage) AddReqToQueue(q *queue.Queue, i interface{}, path string) error {

	return nil
}
