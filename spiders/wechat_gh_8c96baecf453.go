package spiders

/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-28 17:33:09
 * @FilePath: \image_hub\spiders\touxiangshe.go
 * @Description: 微信号自定义处理函数-情侣头像原创榜
 */

//  import (
// 	 "image_hub/model"
// 	 "strings"
//  )

//  // 微信号：touxiangshe, 文章索引: 1 的自定义处理函数
//  // 头像内容处理
//  func seevanlove1(sections []model.Section) []model.Section {

// 	 // Filter out sections with invalid text
// 	 sections = filterDirtyText(sections)

// 	 // 遍历头像当前面的section如果section.Text不为空且含有"头像", section.ImageUrls为空，则将该section.Text赋值为后面的section.ImageUrls不为空的section.Text
// 	 // Concatenate sections containing "头像" with subsequent sections that have non-empty image_urls
// 	 for i := 0; i < len(sections)-1; i++ {
// 		 // Check if the current section has "头像" in its text and has an empty image_urls array
// 		 if len(sections[i].ImageUrls) == 0 && strings.Contains(sections[i].Text, "头像") {
// 			 // Look for the next section that has non-empty image_urls and concatenate the text with a newline separator
// 			 for j := i + 1; j < len(sections); j++ {
// 				 // Only concatenate sections with unique text and don't concatenate the same section more than once
// 				 if sections[j].Text != sections[i].Text && len(sections[j].ImageUrls) > 0 {
// 					 sections[j].Text = sections[i].Text
// 					 break
// 				 }
// 			 }
// 		 }
// 	 }

// 	 // Filter out sections with empty image_urls
// 	 sections = filterEmptyImageUrls(sections)

// 	 // 删掉sections前3项(是网页图标,封面宣传图)
// 	 if len(sections) > 3 {
// 		 sections = sections[3:]
// 	 }

// 	 // 删掉sections最后2项(是网页图标,封面宣传图)
// 	 if len(sections) > 2 {
// 		 sections = sections[:len(sections)-2]
// 	 }

// 	 return sections
//  }
