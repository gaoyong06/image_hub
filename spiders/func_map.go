/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-28 17:33:09
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
	addFunc("touxiangshe1", touxiangshe1)
	addFunc("touxiangshe2", touxiangshe2)
	addFunc("touxiangshe3", touxiangshe3)
	addFunc("touxiangshe4", touxiangshe4)

	// 情侣头像原创榜
	addFunc("seevanlove1", seevanlove1)

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
