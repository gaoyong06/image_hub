/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-03-17 10:16:10
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-21 15:28:45
 * @FilePath: \image_hub\spiders\spider.go
 * @Description: 爬虫接口
 */
package spiders

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var (

	// 已访问的url,避免重复访问
	visited = cmap.New[bool]()
	// make(map[string]bool)
)

type Spider interface {

	// 设置爬虫名称
	SetName(string)

	// 获取爬虫名称
	GetName() string

	// 向队列追求爬取请求
	// q 请求队列
	// e 上级页面HTMLElement,没有时设置为nil
	// baseUrl 请求的基准url,目的是为页面内的相对地址补全为完整的地址
	AddReqToQueue(q *queue.Queue, e *colly.HTMLElement, baseUrl string) error

	// 解析将爬取到的数据至一个规范的结构体中
	// e 当前爬虫请求的返回结果 *colly.HTMLElement 或者  *colly.Response
	// baseUrl 请求的基准url,目的是为页面内的相对地址补全为完整的地址
	ParseData(q *queue.Queue, e interface{}, baseUrl string) (interface{}, error)

	// 业务处理
	// 1. 向队列追加请求
	// 2. 解析数据至结构体
	// 3. 保存数据 或 更新数据 或 继续下一层级的请求
	// e  当前爬虫请求的返回结果 *colly.HTMLElement 或者  *colly.Response
	// baseUrl 请求的基准url,目的是为页面内的相对地址补全为完整的地址
	Process(q *queue.Queue, e interface{}, baseUrl string) error
}
