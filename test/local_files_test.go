/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-24 14:29:12
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-28 11:35:35
 * @FilePath: \image_hub\local_files.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

func TestLocalFiles(i *testing.T) {

	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	panic(err)
	// }

	q, _ := queue.New(
		20, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 1000000}, // Use default queue storage
	)

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c := colly.NewCollector()
	c.WithTransport(t)

	pages := []string{}

	c.OnHTML("h1", func(e *colly.HTMLElement) {
		pages = append(pages, e.Text)
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {

		fmt.Printf("href: %+v\n", e.Attr("href"))
		// c.Visit("file://" + dir + "/html" + e.Attr("href"))
	})

	dir := "D:/work/wechat_download_data/html/test"
	url := "file://" + dir + "/20220526_111900_1.html"
	fmt.Println(url)
	q.AddURL(url)

	// c.Visit(url)
	// c.Wait()

	q.Run(c)

	for i, p := range pages {
		fmt.Printf("%d : %s\n", i, p)
	}
}
