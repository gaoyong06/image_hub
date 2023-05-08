/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-05-08 16:10:39
 * @FilePath: \image_hub\spiders\func_map.go
 * @Description: 爬虫相关公用方法
 */

package spiders

import (
	"fmt"
	"image"
	"image_hub/model"
	"image_hub/pkg/utils"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (

	// 图片类型与理想的尺寸
	imageTypes = map[string]map[string]float64{
		"avatar":     {"width": 400, "height": 400},
		"background": {"width": 1080, "height": 1080},
		"wallpaper":  {"width": 1080, "height": 1920},
		"sticker":    {"width": 300, "height": 300}}

	// 图片尺寸范围
	imageDimensionRange = map[string]map[string]float64{
		"avatar":     {"minWidth": 360, "minHeight": 360, "maxWidth": 1080, "maxHeight": 1200},
		"background": {"minWidth": 500, "minHeight": 500, "maxWidth": 1395, "maxHeight": 1920},
		"wallpaper":  {"minWidth": 600, "minHeight": 900, "maxWidth": 1188, "maxHeight": 2376},
		"sticker":    {"minWidth": 180, "minHeight": 180, "maxWidth": 1080, "maxHeight": 1080},
	}

	// 图片文件大小范围
	imageSizeRange = map[string]map[string]float64{
		"avatar":     {"minSize": 1024 * 20, "maxSize": 1024 * 1024 * 10}, // 20kb~10MB
		"background": {"minSize": 1024 * 20, "maxSize": 1024 * 1024 * 10}, // 20kb~10MB
		"wallpaper":  {"minSize": 1024 * 20, "maxSize": 1024 * 1024 * 20}, // 20kb~20MB
		"sticker":    {"minSize": 1024 * 6, "maxSize": 1024 * 1024 * 4},   // 10kb~4MB
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

// 根据公众号标题和标签,确定公众号文章内图片类型
func GetHtmlImageTypes(htmlStr string) []string {

	// 判断文章内的图片类型都有哪些

	var imageTypes []string

	if strings.Contains(htmlStr, "头像") {
		imageTypes = append(imageTypes, "avatar")
	}

	if strings.Contains(htmlStr, "背景") {
		imageTypes = append(imageTypes, "background", "wallpaper")
	}

	if strings.Contains(htmlStr, "套图") {
		imageTypes = append(imageTypes, "avatar", "background")
	}

	if strings.Contains(htmlStr, "壁纸") {
		imageTypes = append(imageTypes, "wallpaper")
	}

	if strings.Contains(htmlStr, "表情") {
		imageTypes = append(imageTypes, "sticker")
	}

	return imageTypes
}

// 从HTML字符串中解析出Section数组，包含文字和图片
// htmlStr 待解析的html字符串
//
// dataSrcRepeat htmlStr所在的文件目录内的所有html中，重复的图片构成的数组
func ParseSectionsFromHTML(htmlStr string, dataSrcRepeat []string) ([]model.Section, error) {

	// html字符串中图片预期的类型(头像, 背景...)，可选值有avatart, background, wallpaper,sticker
	imageTypes := GetHtmlImageTypes(htmlStr)
	fmt.Printf("================== imageTypes: %#v\n", imageTypes)

	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	var sections []model.Section

	images, err := GetImagesInfoFromHTML(htmlStr)
	if err != nil {
		return nil, err
	}

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

		// 如果当前节点为img标签，提取其中的data-src属性，检查data-src属性值是否在dataSrcRepeat中,如果存在则跳过
		if n.Type == html.ElementNode && n.Data == "img" {

			var dataSrc string
			var src string
			for _, attr := range n.Attr {
				if attr.Key == "data-src" {
					dataSrc = attr.Val
				}

				if attr.Key == "src" {

					src = attr.Val
				}
			}

			// 如果同一个图片在多个网页中重复出现，则可能是宣传图，过滤掉
			if utils.Contains(dataSrcRepeat, dataSrc) {
				fmt.Printf("========================= dataSrcRepeat. src: %s\n", src)
				return true
			}

			// 判断图片的长，宽，大小，宽高比，文件类型(jpg,png,gif...)与图片预期的类型(头像，背景图，壁纸，表情包)是否一致，不一致的过滤掉
			if len(imageTypes) > 0 && !IsValidImage(src, images, imageTypes) {
				fmt.Printf("========================= !IsValidImage. src: %s\n", src)
				return true
			}
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
	return sections, nil
}

// 获取网页内图片的信息，返回一个由图片信息map构成的数组
// map的key如下：
//
//	src: 图片的地址
//	ratio：图片的宽高比
//	width：图片的宽高比
//	height：图片的宽高比
//	format: 图片的格式
//	shape：图片的形状(vertical: 垂直的,horizontal: 水平的,square: 正方形)
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
	var wg sync.WaitGroup
	wg.Add(len(imgTags))
	mutex := &sync.Mutex{}
	for _, imgTag := range imgTags {

		go func(imgTag string) {
			defer wg.Done()

			imgInfo := make(map[string]interface{})

			// 获取图片的源URL
			regex := regexp.MustCompile(`\s+src=["']([^"']*)["']`)
			matches := regex.FindAllStringSubmatch(imgTag, -1)

			if len(matches) < 1 {
				fmt.Println("====== len(srcStr) < 2")
				return
			}

			// 输出匹配到的src属性值
			imgInfo["src"] = matches[0][1]

			// 打开图片文件，读取宽度和高度和大小
			imgFile, err := os.Open(imgInfo["src"].(string))
			if err != nil {
				return
			}
			defer imgFile.Close()
			img, imgFormat, err := image.Decode(imgFile)
			if err != nil {
				return
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

			mutex.Lock()
			imgInfo["ratio"] = imgRatio
			imgInfo["width"] = imgWidth
			imgInfo["height"] = imgHeight
			imgInfo["format"] = imgFormat
			imgInfo["size"] = imgSize
			imgs = append(imgs, imgInfo)
			mutex.Unlock()

		}(imgTag)
	}
	wg.Wait()

	return imgs, nil
}

// 检查图片文件信息是否和imageType一致
func IsValidImage(imgStr string, imagesInfo []map[string]interface{}, imageTypes []string) bool {

	//遍历imagesInfo，找到imagesInfo中，image["src"]==imgStr的项
	for _, image := range imagesInfo {

		if image["src"] == imgStr {

			//得到图片的宽度，高度，宽高比，大小
			imgWidth := image["width"].(float64)
			imgHeight := image["height"].(float64)
			imgSize := image["size"].(float64)
			imgRatio := image["ratio"].(float64)
			imgFormat := image["format"].(string)

			//与上面定义的对应imageTypes其中之一是否匹配，匹配返回true, 全部不匹配则返回false
			for _, imageType := range imageTypes {
				if imageDimensionRange[imageType]["minWidth"] <= imgWidth && imgWidth <= imageDimensionRange[imageType]["maxWidth"] &&
					imageDimensionRange[imageType]["minHeight"] <= imgHeight && imgHeight <= imageDimensionRange[imageType]["maxHeight"] &&
					imageSizeRange[imageType]["minSize"] <= imgSize && imgSize <= imageSizeRange[imageType]["maxSize"] &&

					utils.Contains(imageFormatRange[imageType], imgFormat) &&
					imageRatioRange[imageType]["min"] <= imgRatio && imgRatio <= imageRatioRange[imageType]["max"] {
					return true
				}
			}
			break
		}
	}
	return false

}

// 根据图片的宽度，高度，大小，格式，尺寸，宽高比打分计算推断图片的类型（类型：头像,背景图，壁纸，表情包, 未知）
// 返回一个由图片信息map构成的数组
// map的key如下：
//
//	src: 图片的地址
//	ratio：图片的宽高比
//	width：图片的宽度
//	height：图片的高度
//	format: 图片的格式
//	shape：图片的形状(vertical: 垂直的,horizontal: 水平的,square: 正方形)
//	type：图片的类型(avatar: 头像,background: 背景图,wallpaper: 壁纸,sticker: 表情包, unknown: 未知的类型)
//
// 工作原理：
//
//	根据宽度、高度、比例、物理空间大小检查图片
//	例如头像的图片，更偏向一个正方形，但是不一定绝对是正方形，只是接近于正方形；而背景图，偏向一个横向的长方形，但是宽和高差异也不是特别大；
//	而手机壁纸是竖向的长方形，宽度小，高度高，高度比宽度要高很多；而表情包，尺寸上，一般比头像小，宽高比和头像相差不大，文件物理尺寸上一般比头像小一些
func InferImageTypeByScoreFromHTML(htmlStr string) ([]map[string]interface{}, []map[string]interface{}, error) {

	imgs, err := GetImagesInfoFromHTML(htmlStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get images info from HTML: %v", err)
	}

	fmt.Printf("============================== GetImagesInfoFromHTML len: %d ==========\n", len(imgs))

	var filteredImgs []map[string]interface{}

	// 判断每种类型的得分并找出得分最高的图片类型
	for _, imgInfo := range imgs {
		imgWidth := imgInfo["width"].(float64)
		imgHeight := imgInfo["height"].(float64)
		imgSize := imgInfo["size"].(float64)
		imgRatio := imgInfo["ratio"].(float64)
		imgFormat := imgInfo["format"].(string)

		maxScore, maxType := 0.0, "unknown"

		for typeName, dimRange := range imageDimensionRange {
			if (dimRange["minWidth"] <= imgWidth) && (imgWidth <= dimRange["maxWidth"]) && (dimRange["minHeight"] <= imgHeight) && (imgHeight <= dimRange["maxHeight"]) &&
				(imgSize >= imageSizeRange[typeName]["minSize"]) && (imgSize <= imageSizeRange[typeName]["maxSize"]) &&
				utils.Contains(imageFormatRange[typeName], imgFormat) &&
				(imageRatioRange[typeName]["min"] <= imgRatio) && (imgRatio <= imageRatioRange[typeName]["max"]) {

				// 计算图片得分
				widthScore := (imageTypes[typeName]["width"] - math.Abs(imageTypes[typeName]["width"]-imgWidth)) / imageTypes[typeName]["width"]
				heightScore := (imageTypes[typeName]["height"] - math.Abs(imageTypes[typeName]["height"]-imgHeight)) / imageTypes[typeName]["height"]
				ratioScore := 1.0 - math.Abs(imgRatio-imageRatioRange[typeName]["max"])/imageRatioRange[typeName]["max"]
				ratioDeviation := 1.0 / (1.0 + math.Abs(imgRatio-imageRatioRange[typeName]["min"]))
				score := widthScore + heightScore + ratioScore + ratioDeviation

				// 若得分高于之前图片则设为最高得分
				if score > maxScore {
					maxScore, maxType = score, typeName
				}
			}
		}

		// 生成图片信息map
		imgInfo["type"] = maxType

		// 归类并筛选出符合条件的图片信息
		if (maxType == "unknown") || ((imgInfo["size"].(float64) < imageSizeRange[maxType]["minSize"]) || (imgInfo["size"].(float64) > imageSizeRange[maxType]["maxSize"]) ||
			(imgInfo["ratio"].(float64) < imageRatioRange[maxType]["min"]) || (imgInfo["ratio"].(float64) > imageRatioRange[maxType]["max"])) {
			filteredImgs = append(filteredImgs, imgInfo)
		}
	}

	return imgs, filteredImgs, nil
}

// 根据HTML文本提取所有图片的信息，并返回符合要求的图片信息及被过滤的图片信息
func GetFilteredImagesInfoFromHTML(htmlStr string, expectedType string, allowedImgFormat []string) ([]map[string]interface{}, []map[string]interface{}, error) {

	// 调用GetImagesInfoFromHTML提取所有图片信息
	imgs, err := GetImagesInfoFromHTML(htmlStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get images info from HTML: %v", err)
	}

	// 处理不符合要求的图片，将其从imgs中删除，并添加到filteredImgs
	var filteredImgs []map[string]interface{}
	for i := len(imgs) - 1; i >= 0; i-- {

		img := imgs[i]

		// 检查图片格式是否符合要求。如果不符合要求，将其从imgs中删除
		count := lo.Count(allowedImgFormat, img["format"].(string))
		isAllowedImgFormat := false
		if count > 0 { // not allowed format, remove image from imgs and add to filteredImgs.  (This is a reverse of RemoveImageFrom
			isAllowedImgFormat = true
		}

		// img["type"]只是参考, 还依赖于文件大小与宽高比
		if img["type"] != expectedType && (img["size"].(float64) < imageSizeRange[expectedType]["minSize"] ||
			img["size"].(float64) > imageSizeRange[expectedType]["maxSize"] ||
			img["ratio"].(float64) < imageRatioRange[expectedType]["min"] ||
			img["ratio"].(float64) > imageRatioRange[expectedType]["max"]) || !isAllowedImgFormat {

			filteredImgs = append(filteredImgs, img)
			imgs = append(imgs[:i], imgs[i+1:]...)
		}
	}

	return imgs, filteredImgs, nil
}

// 读取directoryPath所有html的文件将各个文件中的img标签的data-src内的值取出来如果重复出现(出现次数大于1),则记录下来返回
func GetImageDataSrcRepeat(directoryPath string) ([]string, error) {

	// 检查directoryPath是否是目录
	info, err := os.Stat(directoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory info: %v", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("directoryPath argument is not a directory")
	}

	// 记录重复的数量，如果大于1则表示有重复的图片
	maxRepeated := 1

	// First, we need to get a list of all HTML files in the given directory
	htmlFiles, err := filepath.Glob(directoryPath + "*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to get a list of all HTML files: %v", err)
	}

	// Initialize a map to store the data-src values and their occurrences
	dataSrcOccurrences := make(map[string]int)

	// Loop through each HTML file asynchronously
	var wg sync.WaitGroup
	lock := sync.Mutex{}
	for _, file := range htmlFiles {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			// Read the contents of the file
			contents, err := ioutil.ReadFile(file)
			if err != nil {
				log.Errorf("Failed to read contents of file %s: %v", file, err)
				return
			}

			regexStart := `(?i)<img.*?data-src\s*=\s*("|')([^"']+)("|')`
			// Split the HTML file into chunks to avoid overflowing memory
			chunkLength := len(contents) / 4 // Split into 4 chunks
			for i := 0; i < len(contents); i += chunkLength {
				chunkEnd := i + chunkLength
				if chunkEnd > len(contents) {
					chunkEnd = len(contents)
				}
				chunk := contents[i:chunkEnd]
				// Use regex to find all the data-src values in the chunk
				re := regexp.MustCompile(regexStart)
				dataSrcValues := make([][]string, 0)
				// Find all <img> tags in the chunk
				imgTags := regexp.MustCompile(`(?i)<img.*?>`).FindAllString(string(chunk), -1)
				for _, tag := range imgTags {
					// Use regex to find the data-src value in each <img> tag
					match := re.FindStringSubmatch(tag)
					if len(match) == 4 {
						dataSrcValues = append(dataSrcValues, []string{match[1], match[2]})
					}
				}

				lock.Lock()
				for _, match := range dataSrcValues {
					dataSrc := match[1]
					dataSrcOccurrences[dataSrc]++
				}
				lock.Unlock()
			}

		}(file)
	}
	wg.Wait()

	// Find the Repeat of all the data-src sets
	dataSrcRepeat := make([]string, 0)
	for dataSrc, occurrences := range dataSrcOccurrences {
		if occurrences > maxRepeated && !utils.Contains(dataSrcRepeat, dataSrc) {
			dataSrcRepeat = append(dataSrcRepeat, dataSrc)
		}
	}

	return dataSrcRepeat, nil
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
