/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-03-18 11:01:46
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-03-21 15:19:35
 * @FilePath: \car_hub\pkg\utils\image.go
 * @Description: 图片处理工具类
 */
package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func SaveImage(r *colly.Response, imageDir string) error {

	// 图片保存
	contentType := r.Headers.Get("Content-Type")
	if strings.HasPrefix(contentType, "image") {

		// url := "https://2sc2.autoimg.cn/escimg/g25/M00/B0/02/f_900x675_0_q87_autohomecar__ChxkqWP5azqAeKCuAAXyUEWOqxM350.jpg"
		url := r.Request.URL

		// urlPath: /escimg/g24/M01/A7/BD/f_900x675_0_q87_autohomecar__Chxky2P5a0CAeR6VAAWJlvDtIW4760.jpg
		urlPath := url.Path

		// dirPath: \escimg\g24\M01\A7\BD
		dirPath := filepath.Dir(urlPath)

		// fileName: f_900x675_0_q87_autohomecar__Chxky2P5a0CAeR6VAAWJlvDtIW4760.jpg
		fileName := filepath.Base(urlPath)

		// 去掉文件名中的"_autohomecar__"
		fileName = strings.ReplaceAll(fileName, "_autohomecar__", "_autoevol_")

		// fileDir: D:\work\images\escimg\g24\M01\A7\BD
		fileDir := filepath.Join(imageDir, dirPath)

		err := os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			log.Errorf("SaveImage failed. err: %+v\n", err)
			return err
		}

		// filePath: D:\work\images\escimg\g24\M01\A7\BD\f_900x675_0_q87_autohomecar__Chxky2P5a0CAeR6VAAWJlvDtIW4760.jpg
		filePath := filepath.Join(fileDir, fileName)
		r.Save(filePath)
	}
	return nil
}
