/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-03-12 10:01:15
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-03-18 11:56:09
 * @FilePath: \car_hub\spiders\detail_page_test.go
 * @Description: 详情页爬虫单测
 */
package spiders

import (
	"car_hub/pkg/utils"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

func TestSpiderDetailPage(t *testing.T) {

	// 获取可被抓取的域名
	domains := strings.Split(Domains, ",")

	// 全国列表页 Collector
	c := colly.NewCollector(
		colly.AllowedDomains(domains...),
		colly.AllowURLRevisit(),
	)
	c.SetRequestTimeout(120 * time.Second)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       5 * time.Second,
	})

	// create a request queue with 2 consumer threads
	q, _ := queue.New(
		1, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	detailPageSpider := &DetailPageSpider{
		Name: UrlTypeUsedCarDetailPage,
	}

	// baseUrl := "https://www.che168.com"
	// // urlStr := "https://www.che168.com/dealer/85869/46333726.html?pvareaid=108783&pos=5#usercid=500100#userpid=500000"
	// // urlStr := "https://www.che168.com/dealer/499070/47365696.html?pvareaid=108783&pos=1#usercid=510100#userpid=510000"
	urlStr := "https://www.che168.com/dealer/489114/47245676.html?pvareaid=100519&userpid=510000&usercid=0&offertype=650&offertag=0&activitycartype=11#pos=20#page=1#rtype=10#isrecom=0#filter=29#module=10#refreshid=0#recomid=0#queryid=1679016598$B$6e930327-cb2d-4d0d-94b4-eb8674d59c3b$37297#cartype=70"

	c.OnHTML("html", func(e *colly.HTMLElement) {

		urlType := e.Response.Ctx.Get(UrlTypeKey)
		if urlType == UrlTypeUsedCarDetailPage {
			detailPageSpider.Process(q, e, BaseUrl)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {

		fmt.Printf("OnResponse: %s, %d bytes\n", r.Request.URL, len(r.Body))
		// 图片保存
		utils.SaveImage(r, ImageDir)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error", err.Error())
	})

	url, _ := url.Parse(urlStr)
	req := &colly.Request{
		URL:    url,
		Method: "GET",
		Ctx:    colly.NewContext(),
	}

	req.Ctx.Put(UrlTypeKey, UrlTypeUsedCarDetailPage)
	q.AddRequest(req)

	err := q.Run(c)
	if err != nil {
		fmt.Printf("Queue.Run() return an error: %v", err)
	}

	isEmpty := q.IsEmpty()
	size, err := q.Size()
	if err != nil {
		fmt.Printf("Queue.Size() return an error: %v", err)
	}
	threads := q.Threads
	fmt.Printf("================== Done. q.IsEmpty: %+v, q.Size: %d, q.Threads: %d ==================\n", isEmpty, size, threads)
}
