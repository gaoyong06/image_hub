package spiders

import (
	"image_hub/model"
	"strings"
)

var funcMap = make(map[string]func(sections []model.Section) []model.Section)

func init() {

	// 头像社-第1条内容自定义方法
	addFunc("touxiangshe1", touxiangshe1)
	addFunc("touxiangshe2", touxiangshe2)
	addFunc("touxiangshe3", touxiangshe3)
	addFunc("touxiangshe4", touxiangshe4)
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

// --------------------------------- 公用方法 ------------------------------------

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
