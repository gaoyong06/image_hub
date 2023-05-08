package test

import (
	"encoding/json"
	"fmt"
	"image_hub/spiders"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestParseSections(t *testing.T) {

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
	sections, err := spiders.ParseSectionsFromHTML(htmlStr1, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

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
