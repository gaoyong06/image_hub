/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-05 07:56:26
 * @FilePath: \image_hub\spiders\func_map.go
 * @Description: 微信号自定义处理函数map, key: 微信号+文章索引号, val：自定义处理函数
 */

package spiders

import (
	"fmt"
	"image_hub/model"
)

var funcMap = make(map[string]func(article *model.TblArticle, sections []model.Section) []model.Section)

func init() {

	// 默认处理方法
	addFunc("simple", simple)

	// 头像社-第1条内容自定义方法
	addFunc("touxiangshe_1", touxiangshe_1)
	addFunc("touxiangshe_2", touxiangshe_2)
	addFunc("touxiangshe_3", touxiangshe_3)
	addFunc("touxiangshe_4", touxiangshe_4)

	// 情侣头像原创榜
	addFunc("seevanlove_1", seevanlove_1)

	// 头像有点好看
	addFunc("gh_8c96baecf453_1", gh_8c96baecf453)
	addFunc("gh_8c96baecf453_2", gh_8c96baecf453)

	// 头像即新欢
	addFunc("gh_22c17e1db325_1", gh_22c17e1db325)
	addFunc("gh_22c17e1db325_2", gh_22c17e1db325)
	addFunc("gh_22c17e1db325_3", gh_22c17e1db325)
	addFunc("gh_22c17e1db325_4", gh_22c17e1db325)

	// 头像库
	addFunc("touxiangcool_1", touxiangcool)
	addFunc("touxiangcool_2", touxiangcool)
	addFunc("touxiangcool_3", touxiangcool)
	addFunc("touxiangcool_4", touxiangcool_4)

	// 下面的公众号使用sections的默认处理方法simple
	// 头像文案
	// 你的小众头像
	// 换头像bo
	// 每日新头像
	// 要啥头像
	// 琉柒头像
	// 头像娣
	// 女生头像壁纸控
	// 头像先生
	// 头像味
	// 小怪兽头像
	// 二次元头像集

}

// AddFunc adds a custom function to the funcMap
func addFunc(key string, val func(article *model.TblArticle, sections []model.Section) []model.Section) {
	funcMap[key] = val
}

// getFuncKey by wechatId, fileIdx
func getFuncKey(wechatId, fileIdx string) string {

	funcKeyDefault := "simple"
	key := fmt.Sprintf("%s_%s", wechatId, fileIdx)
	_, ok := funcMap[key]
	if !ok {
		key = funcKeyDefault
	}
	return key
}

// RunFunc runs the custom function associated with the given name
func runFunc(key string, article *model.TblArticle, sections []model.Section) []model.Section {

	val, ok := funcMap[key]
	if ok {
		sections = val(article, sections)
	}

	return sections
}
