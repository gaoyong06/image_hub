/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-03-12 16:04:07
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-24 11:16:46
 * @FilePath: \image_hub\pkg\utils\text.go
 * @Description:  文字处理工具类
 */
package utils

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 获取html标签后面文本内容,返回第一个文本内容
func GetNodeTextOne(s *goquery.Selection) string {

	nodeTexts := GetNodeText(s)
	if len(nodeTexts) > 0 {
		return nodeTexts[0]
	}
	return ""

}

// 获取html标签后面文本内容
// 下面获到"[2021年07月 bbb ddd]",s定位到li元素
// <ul class="basic-item-ul" >
//
//		<li><span class="item-name">上牌时间</span>2021年07月</li>
//	 <li><span class="item-name">aaa</span>bbb</li>
//	 <li><span class="item-name">ccc</span>ddd</li>
//		...
//
// </ul>
// https://github.com/PuerkitoBio/goquery/issues/287
func GetNodeText(s *goquery.Selection) []string {

	var nodeTexts []string
	s.Contents().Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "#text" {
			nodeTexts = append(nodeTexts, s.Text())

		}
	})
	return nodeTexts
}

// 去除所有空格
// https://stackoverflow.com/questions/65533097/replace-nbsp-or-0xao-with-space-in-a-string
func RemoveSpace(str string) string {

	newStr := strings.ReplaceAll(str, "\u00a0", "")
	return newStr
}

// 将字符串"¥1万,¥1.6万" 转化为数字"10000,16000"
// 将字符串"1万,1.6万" 转化为数字"10000,16000"
func ConvertTenThousand(str string) (int, error) {

	// 去除所有空格
	newStr := RemoveSpace(str)

	// 去除左边"¥"
	newStr = strings.TrimLeft(newStr, "¥")

	// 去除左边"￥"
	newStr = strings.TrimLeft(newStr, "￥")

	// 去除左边"¥"
	newStr = strings.TrimLeft(newStr, "¥")

	// 去除右边"万"
	newStr = strings.TrimRight(newStr, "万")

	// 转数字乘以10000
	figure, err := strconv.ParseFloat(newStr, 64)
	if err != nil {
		return 0, err
	}
	newFigure := int(figure * 10000)

	return newFigure, nil
}

// 将参考价格范围字符串8.22-10.25 转为 82000,102500
func ConvertTenThousandRanges(str string) (int, int, error) {

	strSlice := strings.Split(str, "-")
	min, err := ConvertTenThousand(strSlice[0])
	if err != nil {
		return 0, 0, err
	}

	max, err := ConvertTenThousand(strSlice[1])
	if err != nil {
		return min, 0, err
	}

	return min, max, nil
}

// 使用正则表达式匹配所有的html标签，并将其替换为空字符串，从而过滤掉所有的html标签
func FilterHTMLTags(str string) string {
	re := regexp.MustCompile(`(?i)<[^>]*>`)
	return re.ReplaceAllString(str, "")
}

// ####任何困难都能把我打倒##########锻炼➕熬夜＝健康的尸体##信女愿一生荤素搭配####################❤️‍🔥 𝑳𝒐𝒗𝒆 𝒎𝒆 𝒆𝒗𝒆𝒓𝒚𝒅𝒂𝒚###########今天除了呼吸什么也没做#####
// 输入是字符串，返回只有文字组成的数组，去掉#号
func GetTexts(str string) []string {
	// 使用正则表达式匹配所有的非#字符，并将其拼接成字符串
	re := regexp.MustCompile(`[^#]+`)
	newStr := strings.Join(re.FindAllString(str, -1), " ")
	// 将字符串按空格分割成数组
	return strings.Split(newStr, " ")
}
