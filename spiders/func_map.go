/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-02 22:13:31
 * @FilePath: \image_hub\spiders\func_map.go
 * @Description: 微信号自定义处理函数map, key: 微信号+文章索引号, val：自定义处理函数
 */

package spiders

import (
	"image_hub/model"
)

var funcMap = make(map[string]func(sections []model.Section) []model.Section)

func init() {

	// 头像社-第1条内容自定义方法
	addFunc("touxiangshe_1", touxiangshe_1)
	addFunc("touxiangshe_2", touxiangshe_2)
	addFunc("touxiangshe_3", touxiangshe_3)
	addFunc("touxiangshe_4", touxiangshe_4)

	// 情侣头像原创榜
	addFunc("seevanlove_1", seevanlove_1)

	//
	addFunc("gh_8c96baecf453_1", gh_8c96baecf453)
	addFunc("gh_8c96baecf453_2", gh_8c96baecf453)

}

// AddFunc adds a custom function to the funcMap
func addFunc(key string, val func(sections []model.Section) []model.Section) {
	funcMap[key] = val
}

// RunFunc runs the custom function associated with the given name
func runFunc(key string, sections []model.Section) []model.Section {

	val, ok := funcMap[key]
	if ok {
		sections = val(sections)
	}

	return sections
}
