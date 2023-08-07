/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date:2023-04-21 18:43:56
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-08-07 17:50:29
 * @FilePath: \image_hub\spiders\base_spider.go
 * @Description: 公众号页面基础爬虫结构体
 */
package spiders

import (
	"fmt"
	"image_hub/model"
	"image_hub/pkg/utils"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	cmap "github.com/orcaman/concurrent-map/v2"
	lop "github.com/samber/lo/parallel"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

var (

	// 已访问的url,避免重复访问
	visited = cmap.New[bool]()
	// make(map[string]bool)

	// tag内的需要被替换为空的特殊字符
	tagDirtyTexts = []string{
		"#",
		"☺︎",
	}

	// 如果文字中含有下面广告关键字则直接跳过,不做处理
	adKeywords = []string{
		"优惠的活动",
		"扫码选礼物",
		"长按扫码即可添加领取",
		"博主朋友圈巨宝藏哦",
		"铂金之恋",
		"绘画学习",
	}

	// section的text内包含下面的文字,则该行文字替换为空字符串
	sectionDirtyTexts = []string{
		// 头像社
		"微信扫一扫关注该公众号",
		"微信号",
		"公众号",
		"长按小图",
		"功能介绍",
		"图源",
		"来自",
		"👇🏻👇🏻👇🏻",
		"@",
		"©️",
		"cr",
		"你们要的",
		"\u200d\u200d",
		"转自",
		"长按保存",
		"点击上方“蓝字”关注我",

		// 情侣头像原创榜
		"情侣头像原创榜",
		"头像即新欢",
		"点击上方蓝色字关注我们",
		"微信",
		"头像研究舍",
		"头像研究舍",
		"点击图片放大，长按图片保存",
		"■",
		"-",
		"。",
		"▼",
		"",
		"男生头像 / 动漫头像 / 壁纸 / 手机壁纸 / 无水印壁纸 / 朋友圈背景图",
		"头像/无水印头像/个性头像 / 明星头像 / 女生头像 /男生头像 / 动漫头像 / 壁纸 / 手机壁纸 / 无水印壁纸 / 朋友圈背景图 头像/无水印头像/个性头像 / 明星头像 / 女生头像 /男生头像 / 动漫头像 / 壁纸 / 手机壁纸 / 无水印壁纸 / 朋友圈背景图",
		"按图片即可保存",
		"01",
		"02",
		"03",
		"04",
		"05",
		"06",
		"07",
		"08",
		"09",
		"10",
		"#",
		"◆",
		"#",
		"75307601",
		"关注我们哦",
		"点击预览",
		"上下滑动信封内纸张",
		"下浏览",
		"点击图片放大",
		"探索粉丝基地精彩内容",
		"一定要收藏",
		"有抽奖",
		"抖音",
		"快手",
	}

	// 微信名和微信号的Map
	nicknameWechatIdMap = map[string]string{

		"头像社":       "touxiangshe",
		"情侣头像原创榜":   "seevanlove",
		"头像有点好看":    "gh_8c96baecf453",
		"头像即新欢":     "gh_22c17e1db325",
		"元気头像":      "NiceWallpaper",
		"头像库":       "touxiangcool",
		"头像文案":      "fashionshijue",
		"你的小众头像":    "h13031h",
		"换头像bo":     "htxb888",
		"每日新头像":     "gh_75640868571b",
		"头像备忘录":     "DNTX9527",
		"小鹿头像酱":     "fairy_goods_thing",
		"梅头像":       "MXLtou",
		"可爱cp头像":    "tcgonglue",
		"要啥头像":      "gh_cdb453299489",
		"情侣头像大全":    "qltxdq",
		"暮昭昭头像馆":    "MzzTxg",
		"琉柒头像":      "lik0894",
		"头像娣":       "Txd777i",
		"精选女生头像":    "touxiang_520",
		"女生头像壁纸控":   "touxiangdiss1",
		"头像先生":      "J79938",
		"头像味":       "gh_bc125df08550",
		"小怪兽头像":     "gh_97a6f9e34972",
		"二次元头像集":    "cpdd52199",
		"头像壁纸大全":    "txbz001",
		"搞怪沙雕头像":    "youtiaotaolu",
		"头像记":       "laixieee",
		"头像酱呀":      "bizhi1994",
		"头像辑":       "touxiangh",
		"头像博主":      "txbozhu",
		"头像号":       "remenyt",
		"百合头像":      "baihetouxiang",
		"背景头像":      "meaijiepai",
		"头像录":       "liaoshangbiji",
		"头像哒":       "gh_367e376abfe0",
		"女生头像宝藏集":   "gh_a600aed1c30d",
		"超火情侣头像":    "chanxuehuiyu",
		"头像壁纸每日推荐":  "touxbizhimeiriTJ",
		"搞怪头像大全":    "gh_089775ff1457",
		"头像微甜":      "txwt-sweet",
		"女生头像壁纸":    "nvshengtouxiang1",
		"ULzzang头像": "Ins-face",
		"头像书":       "Txs5665",
		"头像大全丫":     "wulai969",
		"可爱萌娃头像大全":  "mwtx66695",
		"萌娃头像库":     "gh_bb13ee258433",
		"可爱萌娃头像":    "bqv8897",
		"丸子妹头像":     "bq6691",
		"萌娃表情包可爱":   "bqb598",
		"搞怪头像合集":    "gaoguaitx",
		"沙雕头像君":     "bao_mihuaqi",
		"古风头像控":     "gh_aca0610cf585",
		"古风头像馆":     "gftxg123",
		"古风壁纸馆":     "gfbzg-007",
		"九栀头像":      "afx1990",
		"胖橘子呀":      "Pang-Juziya",
		"玫竹斋":       "meizhuzhai",
		"草莓头像":      "touxiangforever",
	}
)

// 定义公众号页面基础爬虫结构体
// 这里用了面向对象的继承和多态的思想，封装了一个baseSpider
// 后面实现的Spider,就可以拥有相关的方法
// 因为golang不支持虚拟方法(父类调用子类方法),所以在Process方法中,把"子类"的Process,作为第一个参数传进去
// 相关文档
//
//	https://www.codeplayer.org/Wiki/Program/go/%E5%9C%A8Go%E8%AF%AD%E8%A8%80%E9%87%8C%E4%BD%BF%E7%94%A8%E7%BB%A7%E6%89%BF%E7%9A%84%E7%BF%BB%E8%BD%A6%E7%BB%8F%E5%8E%86.html
//	https://hackthology.com/golangzhong-de-mian-xiang-dui-xiang-ji-cheng.html
type baseSpider struct {
	Name string
}

// 设置爬虫名称
func (b *baseSpider) SetName(name string) {
	b.Name = name
}

// 获取爬虫名称
func (b *baseSpider) GetName() string {

	return b.Name
}

// 向队列追求爬取请求
// q 请求队列
// e 上级页面HTMLElement,没有时设置为nil
// baseUrl 请求的基准url,目的是为页面内的相对地址补全为完整的地址
func (b *baseSpider) AddReqToQueue(q *queue.Queue, i interface{}, params map[string]interface{}) error {

	path := params["path"].(string)

	// 目前至支持解析本地文件
	pathUrl := fmt.Sprintf("file://%s", path)

	// 解析 URL
	url, err := url.Parse(pathUrl)
	if err != nil {
		log.Errorf("url.Parse failed. err: %+v\n", err)
		return err
	}

	if _, ok := visited.Get(path); !ok {

		visited.Set(path, true)
		req := &colly.Request{
			URL:    url,
			Method: "GET",
			Ctx:    colly.NewContext(),
		}

		req.Ctx.Put(UrlTypeKey, b.Name)
		q.AddRequest(req)
	}
	return nil
}

// 解析将爬取到的数据至一个规范的结构体中
// e 当前爬虫请求的返回结果 *colly.HTMLElement 或者  *colly.Response
// baseUrl 请求的基准url,目的是为页面内的相对地址补全为完整的地址
func (b *baseSpider) ParseData(q *queue.Queue, i interface{}, params map[string]interface{}) (interface{}, error) {

	// 解析返回html结果
	article := &model.TblArticle{}
	var selector string

	e, ok := i.(*colly.HTMLElement)
	if !ok {
		return nil, fmt.Errorf("invalid type: %T, expected *colly.HTMLElement", i)
	}

	// 文章标题
	selector = "h1#activity-name"
	title := e.ChildText(selector)
	article.Title = title

	// 作者
	selector = "a#js_name"
	author := e.ChildText(selector)
	article.Author = author

	// 微信号
	var wechatId string
	selector = ".profile_meta_value"
	profileMetaValues := e.ChildTexts(selector)
	if len(profileMetaValues) == 0 {
		wechatId = nicknameWechatIdMap[author]
	} else {
		wechatId = profileMetaValues[0]
	}
	article.WechatId = wechatId

	// 收录于合集
	selector = ".article-tag__item"
	tags := e.ChildTexts(selector)

	lop.ForEach(tags, func(tag string, i int) {

		lop.ForEach(tagDirtyTexts, func(text string, j int) {
			tag = strings.ReplaceAll(tag, text, "")
		})
		tags[i] = tag
	})

	article.Tags = tags

	// 发布时间
	publishTime, _ := utils.GetPublishTime(e.Text)
	article.PublishTime = time.Unix(publishTime, 0)

	// fmt.Printf("================ ParseData: url: %s, title: %s\n", url, title)

	// <meta content="http://mp.weixin.qq.com/s?__biz=MjM5NzAyMDIwMA==&amp;mid=2653562471&amp;idx=1&amp;sn=5a209eca9a0c9d92d484dadfa516a807&amp;chksm=bd3ed1208a49583679dddb80f504983511b6bc9d63c89242dd3df68daebd587a78b8fea1afa0#rd"/>
	selector = "meta[property='og:url']"
	ogUrl := e.ChildAttr(selector, "content")
	queryParams, err := utils.GetArticleUrlQueryParams(ogUrl)
	if err != nil {
		log.Errorf("utils.GetArticleUrlQueryParams failed. ogUrl: %s,  err: %+v\n", ogUrl, err)
		return nil, err
	}
	idx := queryParams.Get("idx")
	sn := queryParams.Get("sn")
	biz := queryParams.Get("__biz")
	mid := queryParams.Get("mid")

	article.Idx = cast.ToUint(idx)
	article.Sn = sn
	article.Biz = biz
	article.Mid = cast.ToUint64(mid)
	article.LocalPath = e.Request.URL.String()

	return article, nil
}

// 业务处理
// 1. 向队列追加请求
// 2. 解析数据至结构体
// 3. 保存数据 或 更新数据 或 继续下一层级的请求
// e  当前爬虫请求的返回结果 *colly.HTMLElement 或者  *colly.Response
// baseUrl 请求的基准url,目的是为页面内的相对地址补全为完整的地址
// golang不支持虚拟方法(父类调用子类方法),所以在Process方法中,把"子类"的Process,作为第一个参数传进去
// params 自定义参数,向下层业务传递参数
func (b *baseSpider) Process(s Spider, q *queue.Queue, i interface{}, params map[string]interface{}) error {

	e, ok := i.(*colly.HTMLElement)
	if !ok {
		return fmt.Errorf("%s invalid type: %T, expected *colly.HTMLElement", s.GetName(), i)
	}

	// 解析返回json结果
	article, err := s.ParseData(q, e, params)
	if err != nil {
		log.Errorf("%s ParseData failed. err: %s, url: %+v\n", s.GetName(), err, e.Request.URL.String())
		return err
	}

	if article != nil {

		// 类型断言进行转换
		tblArticle, ok := article.(*model.TblArticle)
		if ok {

			// 如果数据表不存在,则新建数据表
			err := tblArticle.CreateTableIfNotExists()
			if err != nil {

				log.Errorf("%s article.CreateTableIfNotExists failed. err: %s\n", s.GetName(), err)
				fmt.Printf("%s article.CreateTableIfNotExists failed. err: %s\n", s.GetName(), err)
				return err
			}

			// 保存数据
			// 保存到本地article
			sn, err := tblArticle.CreateOrUpdate()
			if err != nil {

				log.Errorf("%s article.CreateOrUpdate failed. err: %s\n", s.GetName(), err)
				fmt.Printf("%s article.CreateOrUpdate failed. err: %s\n", s.GetName(), err)
				return err
			}

			fmt.Printf("%s article.CreateOrUpdate success. sn: %s\n", s.GetName(), sn)
			log.Infof("%s article.CreateOrUpdate success. sn: %s\n", s.GetName(), sn)

			// 按照多个section保存至content_service
			// TODO:调用content_service API完成批量写入

			return nil

		} else {

			fmt.Printf("%s failed to convert article to tblArticle", s.GetName())
			return fmt.Errorf("%s failed to convert article to tblArticle", s.GetName())
		}
	}

	return nil
}
