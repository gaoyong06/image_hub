package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Section struct {
	Text      string
	ImageUrls []string
}

// 从HTML字符串中解析出Section数组，包含文字和图片
func parseSectionsFromHTML(htmlStr string) []Section {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		log.Fatal(err)
	}

	var sections []Section

	// 字符串过滤器，过滤掉不需要的标签，包括空的 span、不可见文本元素等
	filter := func(n *html.Node) bool {
		return n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style" || n.Data == "head" || n.Data == "title" || n.Data == "meta") ||
			n.Type == html.TextNode && strings.TrimSpace(n.Data) == "\u200d"
	}

	var parseNode func(*html.Node, bool)
	parseNode = func(n *html.Node, skip bool) {
		if filter(n) {
			skip = true
		} else if skip {
			return
		} else if n.Type == html.ElementNode && n.Data == "img" {

			// 如果当前节点为img标签，提取其中的src属性作为Section的图片url
			var imageUrl string
			for _, attr := range n.Attr {

				if attr.Key == "src" {

					imageUrl = attr.Val
				}
			}

			// 将图片url添加到当前Section的ImageUrls列表
			if len(sections) <= 0 {
				sections = append(sections, Section{
					Text:      "",
					ImageUrls: []string{},
				})
			}
			currentSection := sections[len(sections)-1]
			currentSection.ImageUrls = append(currentSection.ImageUrls, imageUrl)
			sections[len(sections)-1] = currentSection

		} else if n.Type == html.TextNode && strings.TrimSpace(n.Data) != "" && strings.TrimSpace(n.Data) != "\u200d" {

			// 如果当前节点为文本节点，提取其中的文字内容作为Section的文本内容
			currentText := strings.TrimSpace(n.Data)

			// 创建一个新的Section，并添加到数组中
			newSection := Section{
				Text:      currentText,
				ImageUrls: []string{},
			}
			sections = append(sections, newSection)
		}

		// 递归调用parseNode处理当前节点的所有子节点
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNode(c, skip)
		}
	}

	// 从根节点开始遍历
	parseNode(doc, false)
	return sections
}

func main() {

	url := "http://192.168.1.3/images/20221222_151433_1.html"

	// 发送http GET请求，获取html内容
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	htmlBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	htmlStr1 := string(htmlBytes)

	// 解析HTML字符串为Section数组
	sections := parseSectionsFromHTML(htmlStr1)

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
}
