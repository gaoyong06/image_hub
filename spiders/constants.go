/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-03-12 22:35:53
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-04-21 14:38:30
 * @FilePath: \car_hub\spiders\constants.go
 * @Description: 常量定义
 */
package spiders

const (

	// 域名 英文逗号分隔多个
	Domains string = ""

	// 目标网站域名
	BaseUrl string = ""

	// 图片保存路径
	ImageDir = "D:\\work\\images"
)

// 为了便于在异步请求返回中对不同url做不同的处理，对url做了如下分类
const (

	// 页面类型定义
	UrlTypeKey = "urlType"

	// 二手车列表页
	UrlTypeChinaListPage = "chinaListPage"
)
