/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-04 18:48:06
 * @FilePath: \image_hub\spiders\func_map.go
 * @Description: 微信号自定义处理函数map, key: 微信号+文章索引号, val：自定义处理函数
 */

package spiders

import (
	"image_hub/model"
)

var funcMap = make(map[string]func(article *model.TblArticle, sections []model.Section) []model.Section)

func init() {

	// 头像社-第1条内容自定义方法
	addFunc("touxiangshe_1", touxiangshe_1)
	addFunc("touxiangshe_2", touxiangshe_2)
	addFunc("touxiangshe_3", touxiangshe_3)
	addFunc("touxiangshe_4", touxiangshe_4)

	// 情侣头像原创榜
	addFunc("seevanlove_1", seevanlove_1)

	// 头像有点好看
	addFunc("gh_8c96baecf453_1", gh_8c96baecf453)
	addFunc("gh_8c96baecf453_2", gh_8c96baecf453)

	// 头像即新欢
	addFunc("gh_22c17e1db325_1", gh_22c17e1db325)
	addFunc("gh_22c17e1db325_2", gh_22c17e1db325)
	addFunc("gh_22c17e1db325_3", gh_22c17e1db325)
	addFunc("gh_22c17e1db325_4", gh_22c17e1db325)

	// 头像库
	addFunc("touxiangcool_1", touxiangcool)
	addFunc("touxiangcool_2", touxiangcool)
	addFunc("touxiangcool_3", touxiangcool)
	addFunc("touxiangcool_4", touxiangcool_4)

	// 头像文案
	addFunc("fashionshijue_1", fashionshijue)
	addFunc("fashionshijue_2", fashionshijue)
	addFunc("fashionshijue_3", fashionshijue)
	addFunc("fashionshijue_4", fashionshijue)
	addFunc("fashionshijue_5", fashionshijue)

	// 你的小众头像
	addFunc("h13031h_1", h13031h)
	addFunc("h13031h_2", h13031h)
	addFunc("h13031h_3", h13031h)
	addFunc("h13031h_4", h13031h)

	// 换头像bo
	addFunc("htxb888_1", htxb888)
	addFunc("htxb888_2", htxb888)
	addFunc("htxb888_3", htxb888)
	addFunc("htxb888_4", htxb888)
	addFunc("htxb888_5", htxb888)

	// 每日新头像
	addFunc("gh_75640868571b_1", gh_75640868571b)
	addFunc("gh_75640868571b_2", gh_75640868571b)
	addFunc("gh_75640868571b_3", gh_75640868571b)

	// 要啥头像
	addFunc("gh_cdb453299489_1", gh_cdb453299489)
	addFunc("gh_cdb453299489_2", gh_cdb453299489)
	addFunc("gh_cdb453299489_3", gh_cdb453299489)
	addFunc("gh_cdb453299489_4", gh_cdb453299489)

	// 琉柒头像
	addFunc("lik0894_1", lik0894)

	// 头像娣
	addFunc("Txd777i_1", Txd777i)
	addFunc("Txd777i_2", Txd777i)
	addFunc("Txd777i_3", Txd777i)
	addFunc("Txd777i_4", Txd777i)
	addFunc("Txd777i_5", Txd777i)
	addFunc("Txd777i_6", Txd777i)

	// 女生头像壁纸控
	addFunc("touxiangdiss1_1", touxiangdiss1)
	addFunc("touxiangdiss1_2", touxiangdiss1)
	addFunc("touxiangdiss1_3", touxiangdiss1)
	addFunc("touxiangdiss1_4", touxiangdiss1)
	addFunc("touxiangdiss1_5", touxiangdiss1)

	// 头像先生
	addFunc("J79938_1", J79938)
	addFunc("J79938_2", J79938)
	addFunc("J79938_3", J79938)

	// 头像味
	addFunc("gh_bc125df08550_1", gh_bc125df08550)
	addFunc("gh_bc125df08550_2", gh_bc125df08550)

	// 小怪兽头像
	addFunc("gh_97a6f9e34972_1", gh_97a6f9e34972)
	addFunc("gh_97a6f9e34972_2", gh_97a6f9e34972)
	addFunc("gh_97a6f9e34972_3", gh_97a6f9e34972)
	addFunc("gh_97a6f9e34972_4", gh_97a6f9e34972)
	addFunc("gh_97a6f9e34972_5", gh_97a6f9e34972)
	addFunc("gh_97a6f9e34972_6", gh_97a6f9e34972)
	addFunc("gh_97a6f9e34972_7", gh_97a6f9e34972)
	addFunc("gh_97a6f9e34972_8", gh_97a6f9e34972)

}

// AddFunc adds a custom function to the funcMap
func addFunc(key string, val func(article *model.TblArticle, sections []model.Section) []model.Section) {
	funcMap[key] = val
}

// RunFunc runs the custom function associated with the given name
func runFunc(key string, article *model.TblArticle, sections []model.Section) []model.Section {

	val, ok := funcMap[key]
	if ok {
		sections = val(article, sections)
	}

	return sections
}
