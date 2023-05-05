package spiders

import "image_hub/model"

var funcMap = make(map[string]func(sections []model.Section) []model.Section)

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

func init() {

	// 头像社-第1条内容自定义方法
	addFunc("touxiangshe1", touxiangshe1)
}
