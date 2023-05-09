/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-05-05 10:01:07
 * @FilePath: \image_hub\spiders\first_page.go
 * @Description: 微信公众号第1条内容抓取-头像
 */

package spiders

import (
	"encoding/json"
	"fmt"
	"image_hub/model"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

type onePage struct {
	*baseSpider
}

// NewOnePage
func NewOnePage(name string) Spider {
	return &onePage{
		baseSpider: &baseSpider{
			Name: name,
		},
	}
}

// 解析将爬取到的数据至一个规范的结构体中
// 对于公众号内内容的处理，有下面几个步骤
//  1. 读取被解析html目录下的所有html文件内的图片，得到重复的图片文件列表(这些图片文件基本都是页面图标，配图，宣传图)
//  2. 遍历各个html文件，读取html文本,通过关键字判断内容是哪些(会有多个)类型. 图片类型包括：头像,背景图，壁纸，表情包4种
//  3. 根据2. 得到的图片内容类型，对该html内的所有图片按相应的类型规格做过滤，规格要求包括：宽度范围，高度范围，宽高比范围，文件大小范围，将不合法的图片过滤掉
//  4. 将过滤后的内容通过ParseSectionsFromHTML处理,解析出html的各个sections
//  5. 在各个公众号的自定义函数内，在对3.解析后的sections做微调，过滤掉不规范的section得到最终的sections(有可能该篇公众号内容是壁纸,但是里面有一个宣传图的尺寸等规格，和壁纸是类似的，在3.处未过滤掉)
//  6. 对sections的微调，一般通过section.Text的文本特征(例如：固定的文案)，或section在sections内的index，这个不可控，得具体情况具体分析，各个公众号，一个公众号内不同时间阶段的文章结构会有不同
//  7. 最终得到sections组装到该html解析到的Article结构体中，就完成了从一个page的html字符串至Article结构体的解析
//
// e *colly.HTMLElement 或者  *colly.Response
func (s *onePage) ParseData(q *queue.Queue, i interface{}, params map[string]interface{}) (interface{}, error) {

	dataSrcRepeat := params["dataSrcRepeat"].([]string)
	articleBase, err := s.baseSpider.ParseData(q, i, params)
	if err != nil {
		return nil, fmt.Errorf("invalid type: %T, expected *colly.HTMLElement", i)
	}

	article, ok := articleBase.(*model.TblArticle)
	if !ok {
		fmt.Printf("%s failed to convert article to tblArticle", s.GetName())
		return nil, fmt.Errorf("%s failed to convert article to tblArticle", s.GetName())
	}

	// Type assertion to convert i to *colly.HTMLElement
	e, ok := i.(*colly.HTMLElement)
	if !ok {
		return nil, fmt.Errorf("invalid type: %T, expected *colly.HTMLElement", i)
	}

	// file://D:/work/wechat_download_data/html/test4/20220810_111900_1.html
	url := e.Request.URL.String()

	// 文章标题
	selector := "h1#activity-name"
	title := e.ChildText(selector)

	// 微信号
	var wechatId string
	selector = ".profile_meta_value"
	profileMetaValues := e.ChildTexts(selector)

	if len(profileMetaValues) == 0 {
		panic(fmt.Sprintf("html class .profile_meta_value element is empty. title: %s, url: %s", title, url))
	} else {
		wechatId = profileMetaValues[0]
	}

	fmt.Printf("================== wechatId: %s,  title: %s, url: %s ================", wechatId, title, url)

	// Get the HTML byte slice of the e element
	htmlBytes := e.Response.Body
	htmlStr := string(htmlBytes)

	// Parse the HTML string to extract the sections
	sections, err := ParseSectionsFromHTML(url, htmlStr, dataSrcRepeat)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sections from HTML: %v", err)
	}

	// 调用每个微信号及其内容索引的自定义方法
	fileIdx := getFileName(url)
	funcKey := fmt.Sprintf("%s%s", wechatId, fileIdx)
	sections = runFunc(funcKey, sections)

	// 将sections以json格式打印出来
	sectionsJson, err := json.Marshal(sections)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sections to json: %v", err)
	}
	fmt.Printf("\n\n=================== Sections JSON======================\n\n%s\n", sectionsJson)

	article.Sections = sections
	return article, nil
}
