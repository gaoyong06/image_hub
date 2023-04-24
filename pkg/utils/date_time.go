/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-03-12 16:04:07
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-24 11:17:05
 * @FilePath: \image_hub\pkg\utils\text.go
 * @Description:  日期时间处理工具类
 */

package utils

import "time"

// 用于将类似2023-04-23 11:19的字符串转化为time.Time类型
func StringToTime(str string) (time.Time, error) {
	layout := "2006-01-02 15:04"
	return time.Parse(layout, str)
}
