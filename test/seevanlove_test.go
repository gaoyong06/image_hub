package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image_hub/spiders"
	"io/ioutil"
	"log"
	"net/url"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func TestSeevanlove(t *testing.T) {

	file := "D:/work/wechat_download_data/html/Dump-0422-20-12-37/20230315_222813_1.html"

	// 读取file的内容
	htmlBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	htmlStr := string(htmlBytes)

	// 解析HTML字符串为Section数组
	sections := spiders.ParseSectionsFromHTML(htmlStr)

	// 打印结果
	for _, section := range sections {
		fmt.Printf("%s %#v\n", section.Text, section.ImageUrls)
	}

	// 使用json打印出Section数组
	jsonSection, err := json.Marshal(sections)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=====================================")
	fmt.Println(string(jsonSection))

	onePageSpider := spiders.NewOnePage(spiders.OnePage)

	// Create a reader from the byte slice of HTML content
	htmlReader := bytes.NewReader(htmlBytes)

	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		panic(err)
	}

	selector := "meta[property='og:title']"
	title, isExist := doc.Find(selector).Attr("content")

	fmt.Printf("================================ title: %s\n", title)

	if isExist {

	} else {
		fmt.Printf("no matching content found for file %s", file)
	}

	url, err := url.Parse("file:///" + file)
	if err != nil {
		log.Fatal(err)
	}

	// doc的整个html，赋值给变量e
	e := &colly.HTMLElement{
		Request: &colly.Request{
			URL: url,
		},
		Response: &colly.Response{
			Body: htmlBytes,
		},
		DOM: doc.Find("html"),
	}

	article, err := onePageSpider.ParseData(nil, e, "")
	if err != nil {

		panic(err)
	}

	// 使用json打印出article
	jsonArticle, err := json.Marshal(article)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("================ jsonArticle =====================")
	fmt.Println(string(jsonArticle))

}
