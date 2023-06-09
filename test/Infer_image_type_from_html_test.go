package test

// 测试程序测试spiders.InferImageTypeFromHTML方法，测试程序的逻辑是：
// 1. 读取一个目录下的所有html文件
// 2. 逐个遍历目录下的各个html文件
// 3. 通过spiders.InferImageTypeFromHTML读取到该html内的所有图片信息(imgsInfo)，需要被过滤的图片信息(filteredImgs)
// 4. 在原来的html基础上，对所有的图片增加浮层显示对应图片信息(imgInfo),显示的信息包括
//图片类型：imgInfo["type"] 取值：avatar，background, wallpaper, sticker
//图片宽高比：imgInfo["ratio"]
//图片宽度：imgInfo["width"]
//图片高度：imgInfo["height"]
//图片文件类型：imgInfo["format"] 取值：jpg,png,jpeg,webp
//图片形状：imgInfo["shape"] 取值：vertical,horizontal,square
//图片文件大小: imgInfo["size"]: 单位：字节，显示时转为KB
//用样式中的position、top、left、width、height等属性控制浮层大小和位置
// 5. 把新的html字符串写入一个新文件，新文件的文件名使用: update_{原来该html文件名}

import (
	"fmt"
	"image_hub/spiders"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestInferImageTypeFromHTML(t *testing.T) {

	// 1. 定义要读取的目录路径
	directoryPath := "D:/work/wechat_download_data/html/test5/Dump-0422-20-12-37/"

	// 2. 读取该目录下的所有文件，除了以"update_开头的文件"
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		fmt.Println(err)
	}

	// 1. 遍历目录中的所有文件，除了以"update_"开头的文件
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "update_") {
			continue
		}

		// 组合出文件的路径
		filePath := filepath.Join(directoryPath, file.Name())

		// 2. 读取文件中的HTML内容
		htmlData, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 将htmlData转换为字符串
		htmlStr := string(htmlData)

		//获取HTML中的图片信息imgsInfo和需要被过滤的图片filteredImgs
		imgsInfo, filteredImgs, err := spiders.InferImageTypeFromHTML(filePath, htmlStr)
		if err != nil {
			t.Error(err)
			continue
		}

		// 4. 在html中添加浮层显示图片信息
		newHtmlStr := addImageInfoOverlayToHTML(htmlStr, imgsInfo, filteredImgs)

		// 5. 把新的html字符串写入一个新文件，新文件的文件名使用: update_{原来该html文件名}
		newFilePath := filepath.Join(directoryPath, fmt.Sprintf("update_%s", file.Name()))
		err = ioutil.WriteFile(newFilePath, []byte(newHtmlStr), 0644)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

// 获取HTML中的图片信息imgsInfo和需要被过滤的图片filteredImgs
func addImageInfoOverlayToHTML(htmlStr string, imgsInfo []map[string]interface{}, filteredImgs map[string]map[string]interface{}) string {

	newHtmlStr := htmlStr
	// 用正则表达式在HTML字符串中查找img标签
	imgRegex, _ := regexp.Compile(`<img.*?src=["|'](.*?)["|'].*?>`)
	imgTags := imgRegex.FindAllString(htmlStr, -1)

	for _, imgTag := range imgTags {

		// 获取当前img标签的内容

		var imgSrc string
		// 获取图片的源URL
		regex := regexp.MustCompile(`\s+src=["']([^"']*)["']`)
		matches := regex.FindAllStringSubmatch(imgTag, -1)

		if len(matches) < 1 {
			fmt.Println("====== len(srcStr) < 2")
			continue
		}

		// 输出匹配到的src属性值
		imgSrc = matches[0][1]

		// 遍历imgsInfo列表，找到对应的图片信息，包括该图片的类型、宽高比、大小、等等
		var curImgInfo map[string]interface{} // 当前img标签匹配到的图片
		for _, img := range imgsInfo {
			if img["src"].(string) == imgSrc {
				curImgInfo = img
				break
			}
		}
		if curImgInfo == nil {
			continue
		}

		// 获取图片信息
		imgSize := curImgInfo["size"].(float64)

		// 组装图片信息叠加浮层的样式
		var borderColor string
		var backgroundColor string

		if _, ok := filteredImgs[imgSrc]; ok {
			borderColor = "red"
			backgroundColor = "rgba(255,0,0,0.5)"
		} else {
			borderColor = "green"
			backgroundColor = "rgba(0,255,0,0.5)"
		}

		// Generate overlay text displaying image information in the top right corner
		imgInfoOverlay := fmt.Sprintf(`
			<div style="position:absolute; top: 0; right: 0; transform: translate(-10px, 10px); padding: 10px; background-color: %s; color: white; line-height:20px !important; font-size: 12px; z-index:999">
		
				<span style="display: block;">Ratio:  %f</span>
				<span style="display: block;">Width:  %f</span>
				<span style="display: block;">Height:  %f</span>
				<span style="display: block;">Format: %s</span>
				<span style="display: block; border: 1px solid white; padding: 2px;">Type: %s</span>
				<span style="display: block;">Shape:  %s</span>
				<span style="display: block;">Size: %s </span>

			</div>`, backgroundColor, curImgInfo["ratio"], curImgInfo["width"], curImgInfo["height"], curImgInfo["format"], curImgInfo["type"], curImgInfo["shape"], convert2KB(imgSize))

		//为新的img标签添加在img标签内的内容、以及img标签自身的class等样式
		newImgTag := fmt.Sprintf(`
			<div style="position:relative">
					<img src="%s" class="%s" style="%s border:2px solid %s;">
				%s
			</div>`,
			imgSrc,
			strings.Join(getClasses(imgTag), " "),
			getStyle(imgTag),
			borderColor,
			imgInfoOverlay,
		)

		// 修改对应img标签的内容
		newHtmlStr = strings.Replace(newHtmlStr, imgTag, newImgTag, -1)
	}

	return newHtmlStr
}

// 获取class属性值
func getClasses(input string) []string {
	classStrings := []string{}
	re := regexp.MustCompile(`class="([^"]*)"`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return classStrings
	}

	classes := matches[1]
	classStrings = strings.Split(classes, " ")

	return classStrings
}

// 获取style属性值
func getStyle(input string) string {
	re := regexp.MustCompile(`style="([^"]*)"`)
	match := re.FindStringSubmatch(input)
	if len(match) < 2 {
		return ""
	}

	return match[1]
}

func convert2KB(size float64) string {
	return fmt.Sprintf("%.2fKB", size/1000.0)
}
