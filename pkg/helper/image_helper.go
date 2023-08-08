/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-05-16 14:40:45
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-08 12:23:08
 * @FilePath: \content_service\pkg\helper\image_helper.go
 * @Description: 图片工具包
 */
package helper

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// relativeImageLocalPath 获取localPath的相对地址，如果localPath是一个windows的本地文件路径，包含\, 或者\\, 将目录分隔符统一替换为/
func GetRelativeImagePath(rootDir string, localPath string) string {
	if strings.ContainsAny(localPath, "\\/") {
		localPath = filepath.ToSlash(localPath)
		localPath = strings.TrimPrefix(localPath, filepath.ToSlash(rootDir))
	}
	return localPath
}

// 获取格式化后的本地地址
func GetFmtLocalPath(localPath string) string {
	if strings.ContainsAny(localPath, "\\/") {
		localPath = filepath.ToSlash(localPath)
	}
	return localPath
}

// GenerateHTML在dirName目录下,生成html,html内都是imageURLs的图
func GenerateHTML(dirName string, imageURLs []string) {
	// 检查目录是否存在
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		// 目录不存在，则创建目录
		err := os.Mkdir(dirName, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 生成HTML文件路径和文件名
	timestamp := time.Now().Format("20060102150405")
	htmlFileName := fmt.Sprintf("index-%s.html", timestamp)
	htmlPath := filepath.Join(dirName, htmlFileName)

	// 生成HTML内容
	htmlContent := "<html><body>\n<div style=\"display: flex; flex-wrap: wrap; justify-content: flex-start; align-items: flex-start;\">"
	for _, imageURL := range imageURLs {
		htmlContent += fmt.Sprintf("<div style=\"flex: 1 0 30%%; margin: 5px;\"><img src=\"%s\" style=\"max-width: 100%%; height: auto;\"></div>\n", imageURL)
	}
	htmlContent += "</div></body></html>"

	// 将HTML内容写入文件
	err := ioutil.WriteFile(htmlPath, []byte(htmlContent), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HTML file generated at: %s\n", htmlPath)
}
