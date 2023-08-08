/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-03-17 10:16:10
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-08 11:22:57
 * @FilePath: \image_hub\spiders\spider.go
 * @Description: 爬虫接口
 */
package spiders

import (
	"github.com/gocolly/colly/v2/queue"
)

type Spider interface {

	// 设置爬虫名称
	SetName(string)

	// 获取爬虫名称
	GetName() string

	// 向队列追求爬取请求
	// q 请求队列
	// e 上级页面HTMLElement,没有时设置为nil
	// extra 自定义参数
	AddReqToQueue(q *queue.Queue, i interface{}, extra map[string]interface{}) error

	// 解析将爬取到的数据至一个规范的结构体中
	// e 当前爬虫请求的返回结果 *colly.HTMLElement 或者  *colly.Response
	// extra 自定义参数
	ParseData(q *queue.Queue, i interface{}, extra map[string]interface{}) (interface{}, error)

	// 业务处理
	// 1. 向队列追加请求
	// 2. 解析数据至结构体
	// 3. 保存数据 或 更新数据 或 继续下一层级的请求
	// e  当前爬虫请求的返回结果 *colly.HTMLElement 或者  *colly.Response
	// extra 自定义参数,向下层业务传递参数
	Process(s Spider, q *queue.Queue, i interface{}, extra map[string]interface{}) error
}
