package spiders

import "image_hub/model"

// touxiangshe1 去掉sections中最后一个section中的imageUrls中的最后一项，然后返回修改后的sections
func touxiangshe1(sections []model.Section) []model.Section {

	if len(sections) > 0 {
		lastSection := &sections[len(sections)-1]
		if len(lastSection.ImageUrls) > 0 {
			lastSection.ImageUrls = lastSection.ImageUrls[:len(lastSection.ImageUrls)-1]
		}
	}
	return sections
}
