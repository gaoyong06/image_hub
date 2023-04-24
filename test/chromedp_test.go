/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-24 16:51:42
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-24 17:38:08
 * @FilePath: \image_hub\test\chromedp_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
)

func TestScrape(t *testing.T) {

	fmt.Println("================================================ TestScrape")

	// Create a new headless Chrome browser instance
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Create a new collector
	c := colly.NewCollector()

	// Set the collector's user agent
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	// Set the collector's request timeout
	c.SetRequestTimeout(30 * time.Second)

	// Set the collector's error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %v\nError: %v\n", r.Request.URL, err)
	})

	// Set the collector's response handler
	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("Visited: %v\n", r.Request.URL)
	})

	// Set the collector's scrape handler
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
	})

	// Add URLs to the collector's queue
	for _, url := range []string{"http://www.baidu.com"} {
		c.Visit(url)
	}

	// Wait for all requests to finish
	c.Wait()

	// Block until all network requests finish
	chromedp.Run(ctx, network.Enable())
}
