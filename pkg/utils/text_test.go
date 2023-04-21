package utils

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var (
	htmlString = ` <ul class="basic-item-ul" >
		<li><span class="item-name">上牌时间</span>2021年07月</li>
		<li><span class="item-name">aaa</span>bbb</li>
		<li><span class="item-name">ccc</span>ddd</li>
	</ul>`
)

func TestGetNodeText(t *testing.T) {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlString))
	if err != nil {
		log.Fatal(err)
	}

	text := GetNodeText(doc.Find(".basic-item-ul > li").Eq(0))
	fmt.Printf("text: %+v\n", text)
}

func TestGetNodeTextOne(t *testing.T) {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlString))
	if err != nil {
		log.Fatal(err)
	}

	text := GetNodeTextOne(doc.Find(".basic-item-ul > li").Eq(0))
	fmt.Printf("text: %+v\n", text)

}
