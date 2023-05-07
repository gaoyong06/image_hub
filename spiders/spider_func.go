/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-28 17:33:09
 * @FilePath: \image_hub\spiders\func_map.go
 * @Description: 爬虫相关公用方法
 */

package spiders

import (
	"fmt"
	"image"
	"image_hub/model"
	"log"
	"math"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (

	// 图片类型与理想的尺寸
	imageTypes = map[string]map[string]int{
		"avatar":     {"width": 400, "height": 400},
		"background": {"width": 1080, "height": 1080},
		"wallpaper":  {"width": 1080, "height": 1920},
		"sticker":    {"width": 300, "height": 300}}

	// 图片尺寸范围
	imageDimensionRange = map[string]map[string]float64{
		"avatar":     {"minWidth": 360, "minHeight": 360, "maxWidth": 1080, "maxHeight": 1200},
		"background": {"minWidth": 500, "minHeight": 500, "maxWidth": 1395, "maxHeight": 1920},
		"wallpaper":  {"minWidth": 864, "minHeight": 1728, "maxWidth": 1188, "maxHeight": 2376},
		"sticker":    {"minWidth": 180, "minHeight": 180, "maxWidth": 1080, "maxHeight": 1080},
	}

	// 图片文件大小范围
	imageSizeRange = map[string]map[string]float64{
		"avatar":     {"minSize": 1024 * 20, "maxSize": 1024 * 1024 * 2}, // 20kb~2MB
		"background": {"minSize": 1024 * 20, "maxSize": 1024 * 1024 * 2}, // 20kb~2MB
		"wallpaper":  {"minSize": 1024 * 20, "maxSize": 1024 * 1024 * 4}, // 20kb~4MB
		"sticker":    {"minSize": 1024 * 6, "maxSize": 1024 * 1024 * 2},  // 10kb~2MB
	}

	// 文件类型范围
	imageFormatRange = map[string][]string{
		"avatar":     {"jpg", "jpeg", "png", "webp"},
		"background": {"jpg", "jpeg", "png", "webp"},
		"wallpaper":  {"jpg", "jpeg", "png", "webp"},
		"sticker":    {"jpg", "jpeg", "png", "webp", "gif"},
	}

	// 图片宽高比范围
	imageRatioRange = map[string]map[string]float64{
		"avatar":     {"min": 0.92, "max": 1.30},
		"background": {"min": 0.80, "max": 1.20},
		"wallpaper":  {"min": 0.97, "max": 2.17},
		"sticker":    {"min": 0.15, "max": 1.14},
	}
)

// 从HTML字符串中解析出Section数组，包含文字和图片
func ParseSectionsFromHTML(htmlStr string) []model.Section {

	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		log.Fatal(err)
	}

	var sections []model.Section

	// 字符串过滤器，过滤掉不需要的标签，包括空的 span、不可见文本元素等, #activity-name，#meta_content，#js_tags 三个标签的过滤
	filter := func(n *html.Node) bool {

		if n.Type == html.ElementNode && n.Data == "script" {
			return true
		}
		if n.Type == html.ElementNode && n.Data == "style" {
			return true
		}
		if n.Type == html.ElementNode && n.Data == "head" {
			return true
		}
		if n.Type == html.ElementNode && n.Data == "title" {
			return true
		}
		if n.Type == html.ElementNode && n.Data == "meta" {
			return true
		}

		if n.Type == html.ElementNode && len(n.Attr) > 0 {
			for _, attr := range n.Attr {
				if attr.Key == "id" && (attr.Val == "activity-name" || attr.Val == "meta_content" || attr.Val == "js_tags") {
					return true
				}
			}
		}

		if n.Type == html.TextNode && strings.TrimSpace(n.Data) == "\u200d" {
			return true
		}
		return false
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
				sections = append(sections, model.Section{
					Text:      "",
					ImageUrls: []string{},
				})
			}

			// imageUrl不为空则追加到ImageUrls中
			if len(imageUrl) > 0 {
				currentSection := sections[len(sections)-1]
				currentSection.ImageUrls = append(currentSection.ImageUrls, imageUrl)
				sections[len(sections)-1] = currentSection
			}

		} else if n.Type == html.TextNode && strings.TrimSpace(n.Data) != "" && strings.TrimSpace(n.Data) != "\u200d" {

			// 如果当前节点为文本节点，提取其中的文字内容作为Section的文本内容
			currentText := strings.TrimSpace(n.Data)

			// 创建一个新的Section，并添加到数组中
			newSection := model.Section{
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

// 获取网页内图片的信息，返回一个由图片信息map构成的数组
// map的key如下：
//
//	ratio：图片的宽高比
//	width：图片的宽高比
//	height：图片的宽高比
//	format: 图片的格式
//	type：图片的类型(avatar: 头像,background: 背景图,wallpaper: 壁纸,sticker: 表情包, unknown: 未知的类型)
//	shape：图片的形状(vertical: 垂直的,horizontal: 水平的,square: 正方形)
//
// 工作原理：
//
//	根据宽度、高度、比例、物理空间大小检查图片，不符合头像、背景图、壁纸、表情包尺寸的图片会被过滤出来
//	例如头像的图片，更偏向一个正方形，但是不一定绝对是正方形，只是接近于正方形；而背景图，偏向一个横向的长方形，但是宽和高差异也不是特别大；
//	而手机壁纸是竖向的长方形，宽度小，高度高，高度比宽度要高很多；而表情包，尺寸上，一般比头像小，宽高比和头像相差不大，文件物理尺寸上一般比头像小一些
//	返回结果是每个图片一个map
func GetImagesInfoFromHTML(htmlStr string) ([]map[string]interface{}, error) {

	// 用正则表达式在HTML字符串中查找img标签
	imgRegex, err := regexp.Compile(`<img.*?src=["|'](.*?)["|'].*?>`)
	if err != nil {
		return nil, fmt.Errorf("failed to compile imgRegex: %v", err)
	}
	imgTags := imgRegex.FindAllString(htmlStr, -1)

	// 遍历每个img标签，通过宽度、高度、比例、物理空间大小过滤图片
	var imgs []map[string]interface{}
	for _, imgTag := range imgTags {

		imgInfo := make(map[string]interface{})

		// 获取图片的源URL
		regex := regexp.MustCompile(`\s+src=["']([^"']*)["']`)
		matches := regex.FindAllStringSubmatch(imgTag, -1)

		if len(matches) < 1 {
			fmt.Println("====== len(srcStr) < 2")
			continue
		}

		// 输出匹配到的src属性值
		imgInfo["src"] = matches[0][1]

		// 打开图片文件，读取宽度和高度和大小
		imgFile, err := os.Open(imgInfo["src"].(string))
		if err != nil {
			continue
		}
		defer imgFile.Close()
		img, imgFormat, err := image.Decode(imgFile)
		if err != nil {

			continue
		}
		imgWidth := float64(img.Bounds().Max.X)  // 获取图片宽度
		imgHeight := float64(img.Bounds().Max.Y) // 获取图片高度
		imgSizeInfo, _ := imgFile.Stat()
		imgSize := float64(imgSizeInfo.Size())

		if imgHeight > imgWidth {
			imgInfo["shape"] = "vertical"
		} else if imgHeight < imgWidth {
			imgInfo["shape"] = "horizontal"
		} else {
			imgInfo["shape"] = "square"
		}

		var imgRatio float64 = 0
		if imgWidth > 0 {
			imgRatio = float64(imgHeight) / float64(imgWidth)
		}

		imgInfo["ratio"] = imgRatio
		imgInfo["width"] = imgWidth
		imgInfo["height"] = imgHeight
		imgInfo["format"] = imgFormat
		imgInfo["size"] = imgSize

		// 判断每种类型的得分
		scores := make(map[string]float64)
		for typeName, dimRange := range imageDimensionRange {
			if dimRange["minWidth"] <= float64(imgWidth) && float64(imgWidth) <= dimRange["maxWidth"] &&
				dimRange["minHeight"] <= float64(imgHeight) && float64(imgHeight) <= dimRange["maxHeight"] &&
				float64(imgSize) >= imageSizeRange[typeName]["minSize"] && float64(imgSize) <= imageSizeRange[typeName]["maxSize"] &&
				func(imageFormat string, formats []string) bool {
					for _, format := range formats {
						if imageFormat == format {
							return true
						}
					}
					return false
				}(imgFormat, imageFormatRange[typeName]) &&
				imageRatioRange[typeName]["min"] <= imgRatio && imgRatio <= imageRatioRange[typeName]["max"] {
				scores[typeName] = (float64(imageTypes[typeName]["width"])-math.Abs(float64(imageTypes[typeName]["width"])-float64(imgWidth)))/float64(imageTypes[typeName]["width"]) +
					(float64(imageTypes[typeName]["height"])-math.Abs(float64(imageTypes[typeName]["height"])-float64(imgHeight)))/float64(imageTypes[typeName]["height"]) +
					1/(1+math.Abs(imgRatio-float64(imageRatioRange[typeName]["min"]))) +
					1/(1+math.Abs(imgRatio-float64(imageRatioRange[typeName]["max"])))
			}
		}

		var maxScore float64
		var maxType string
		for typeName, score := range scores {
			if score > maxScore {
				maxScore = score
				maxType = typeName
			}
		}

		if maxType == "" {
			maxType = "unknown"
		}

		imgInfo["type"] = maxType
		imgs = append(imgs, imgInfo)
	}

	return imgs, nil
}

// 根据HTML文本提取所有图片的信息，并返回符合要求的图片信息及被过滤的图片信息
func GetFilteredImagesInfoFromHTML(htmlStr string, expectedType string) ([]map[string]interface{}, []map[string]interface{}, error) {

	// 调用GetImagesInfoFromHTML提取所有图片信息
	imgs, err := GetImagesInfoFromHTML(htmlStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get images info from HTML: %v", err)
	}

	// 处理不符合要求的图片，将其从imgs中删除，并添加到filteredImgs
	var filteredImgs []map[string]interface{}
	for i := len(imgs) - 1; i >= 0; i-- {
		img := imgs[i]

		// img["type"]只是参考, 还依赖于文件大小与宽高比
		if img["type"] != expectedType && (img["size"].(float64) < imageSizeRange[expectedType]["minSize"] ||
			img["size"].(float64) > imageSizeRange[expectedType]["maxSize"] ||
			img["ratio"].(float64) < imageRatioRange[expectedType]["min"] ||
			img["ratio"].(float64) > imageRatioRange[expectedType]["max"]) {

			filteredImgs = append(filteredImgs, img)
			imgs = append(imgs[:i], imgs[i+1:]...)
		}
	}

	return imgs, filteredImgs, nil
}

//------------------------------------ 私有方法 -------------------------------------------------

// 获取文件名最后的数字
// file://D:/work/wechat_download_data/html/test4/20220810_111900_1.html
func getFileName(filePath string) string {

	// 将文件路径按照"/"分割成数组
	arr := strings.Split(filePath, "/")
	// 获取数组最后一个元素
	last := arr[len(arr)-1]
	// 将最后一个元素按照"."分割成数组
	arr2 := strings.Split(last, ".")
	// 获取数组第一个元素
	fileName := arr2[0]
	// 将文件名最后的数字提取出来
	lastNum := fileName[len(fileName)-1:]
	return lastNum
}

// 过滤sections中的敏感字符串、
// 将含有敏感字符串的section.Text设置为空字符串
func filterDirtyText(sections []model.Section) []model.Section {

	// 过滤字符串
	if len(sections) > 0 {
		for i := len(sections) - 1; i >= 0; i-- {
			if len(sections[i].Text) > 0 {
				for _, dirtyText := range sectionDirtyTexts {
					if strings.Contains(sections[i].Text, dirtyText) {
						sections[i].Text = ""
						break
					}
				}
			}
		}
	}

	return sections
}

// 过滤sections中的section.ImageUrls
// 将sections中section.ImageUrls为空数组的section从sections中剔除
func filterEmptyImageUrls(sections []model.Section) []model.Section {

	// Filter out sections with empty image_urls
	filteredSections := make([]model.Section, 0, len(sections))
	for _, section := range sections {
		if len(section.ImageUrls) > 0 {
			filteredSections = append(filteredSections, section)
		}
	}

	return filteredSections
}
