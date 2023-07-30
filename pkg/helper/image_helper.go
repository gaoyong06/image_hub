/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-05-16 14:40:45
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-06-05 22:59:48
 * @FilePath: \content_service\pkg\helper\image_helper.go
 * @Description: 图片工具包
 */
package helper

import (
	"path/filepath"
	"strings"
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
