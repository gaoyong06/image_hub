package utils

import (
	"fmt"
	"regexp"
)

// 用正则表达式在HTML字符串中查找img标签
func GetImgTagsFromHTML(htmlStr string) ([]string, error) {

	imgRegex, err := regexp.Compile(`<\s*img[^>]*src\s*=\s*["']?([^"']+)["']?[^>]*>`)
	if err != nil {
		return nil, fmt.Errorf("failed to compile imgRegex: %v", err)
	}
	imgTags := imgRegex.FindAllString(htmlStr, -1)

	return imgTags, nil
}
