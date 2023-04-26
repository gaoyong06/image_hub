/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-04-21 18:19:58
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-26 10:37:57
 * @FilePath: \image_hub\spiders\constants.go
 * @Description: 常量定义
 */
package spiders

const (

	// 域名 英文逗号分隔多个
	Domains string = ""

	// 目标网站域名
	BaseUrl string = ""
)

// 为了便于在异步请求返回中对不同url做不同的处理，对url做了如下分类
const (

	// 页面类型定义
	UrlTypeKey = "urlType"

	// 公众号第1条
	FirstPage = "firstPage"

	// 公众号第2条
	SecondPage = "secondPage"

	// 公众号第3条
	ThirdPage = "thirdPage"

	// 公众号第4条
	FourPage = "fourPage"

	// 不处理的页面
	UnknownPage = "unknownPage"
)
