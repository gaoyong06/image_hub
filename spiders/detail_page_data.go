/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-03-13 10:43:19
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-03-21 15:50:47
 * @FilePath: \car_hub\spiders\car_info_clear.go
 * @Description: 二手车详情页汽车信息数据清理
 */
package spiders

import (
	"car_hub/model"
	"car_hub/pkg/utils"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

// 获取二手车标签
func GetCarTags(detailPageHtmlElement *colly.HTMLElement) []model.TblUsedCarTag {

	var carTags []model.TblUsedCarTag
	selI := detailPageHtmlElement.DOM.Find("div.car-tags.tags > .tp-tags.tags-light")
	selI.Each(func(i int, s *goquery.Selection) {

		var carTag model.TblUsedCarTag
		carTag.TagName = utils.GetNodeTextOne(s)
		carTag.TagContent = s.Find(".tag-content").Text()
		carTags = append(carTags, carTag)
	})

	return carTags
}

// 二手车id
func GetCarId(detailPageHtmlElement *colly.HTMLElement) int {

	carIdStr := detailPageHtmlElement.ChildAttr("#car_infoid", "value")
	carId := cast.ToInt(carIdStr)
	return carId
}

// 二手车名称
func GetName(detailPageHtmlElement *colly.HTMLElement) string {

	name := detailPageHtmlElement.ChildAttr("#car_carname", "value")
	return name
}

// 二手车价格
func GetPrice(detailPageHtmlElement *colly.HTMLElement) int64 {

	priceStr := detailPageHtmlElement.ChildAttr("#car_price", "value")
	priceNum, _ := decimal.NewFromString(priceStr)
	priceRet := priceNum.Mul(decimal.NewFromFloat(10000))
	price := priceRet.IntPart()
	return price
}

// 是否包含过户费
func GetIsIncludeTransferFee(detailPageHtmlElement *colly.HTMLElement) string {

	value := detailPageHtmlElement.ChildAttr(".price-transfer.ndy", "value")
	return value
}

// 二手车所属省id
func GetProvinceIdFromDetailPage(detailPageHtmlElement *colly.HTMLElement) int {

	provinceIdStr := detailPageHtmlElement.ChildAttr("#car_pid", "value")
	provinceId := cast.ToInt(provinceIdStr)
	return provinceId
}

// 二手车所属城市id
func GetCityId(detailPageHtmlElement *colly.HTMLElement) int {

	cityIdStr := detailPageHtmlElement.ChildAttr("#car_cid", "value")
	cityId := cast.ToInt(cityIdStr)
	return cityId
}

// 二手车品牌id
func GetMakeId(detailPageHtmlElement *colly.HTMLElement) int {

	makeIdStr := detailPageHtmlElement.ChildAttr("#car_brandid", "value")
	makeId := cast.ToInt(makeIdStr)
	return makeId
}

// 二手车车系id
func GetModelId(detailPageHtmlElement *colly.HTMLElement) int {

	modelIdStr := detailPageHtmlElement.ChildAttr("#car_seriesid", "value")
	modelId := cast.ToInt(modelIdStr)
	return modelId
}

// 二手车车型id
func GetModelDetailIdFromDetailPage(detailPageHtmlElement *colly.HTMLElement) int {

	modelDetailIdStr := detailPageHtmlElement.ChildAttr("#car_specid", "value")
	modelDetailId := cast.ToInt(modelDetailIdStr)
	return modelDetailId
}

// 获取表显里程
func getMileage(detailPageHtmlElement *colly.HTMLElement) int {

	// 表显里程
	mileageStr := detailPageHtmlElement.ChildAttr("#car_mileage", "value")
	mileage, err := utils.ConvertTenThousand(mileageStr)
	if err != nil {
		log.Errorf("getMileage failed. err: %+v\n", err)
	}

	return mileage
}

// 获取签名后的VIN码
func GetVinCodeSign(detailPageHtmlElement *colly.HTMLElement) string {

	vinCodeSign := detailPageHtmlElement.ChildAttr("#hidvincodesign", "value")
	return vinCodeSign
}

// 获取上牌日期
func GetFirstRegDate(detailPageHtmlElement *colly.HTMLElement) string {

	firstRegDate := detailPageHtmlElement.ChildAttr("#car_firstregtime", "value")
	return firstRegDate
}

// 获取过户次数
func GetTransferOwnershipCount(detailPageHtmlElement *colly.HTMLElement) int {

	// 定位到ul元素
	selUl := detailPageHtmlElement.DOM.Find(".all-basic-content.fn-clear").Find(".basic-item-ul")
	selLi := selUl.Eq(1).Find("li").Eq(4)
	text := utils.GetNodeTextOne(selLi)
	title := utils.RemoveSpace(selLi.Find(".item-name").Text())
	str := ""
	if title == "过户次数" {
		str = text
	}

	//  定义正则表达式
	reg := regexp.MustCompile(`\d+`)

	//  提取数字
	numStr := reg.FindString(str)

	//  将字符串转换为整数
	num, _ := strconv.Atoi(numStr)
	return num
}

// 二手车卖家id
func GetDealerId(detailPageHtmlElement *colly.HTMLElement) int {

	// 二手车卖家id
	dealerIdStr := detailPageHtmlElement.ChildAttr("#car_dealerid", "value")
	dealerId := cast.ToInt(dealerIdStr)
	return dealerId
}

// 注册用户id
func GetMemberId(detailPageHtmlElement *colly.HTMLElement) int {

	var memberId int
	memberIdStr := detailPageHtmlElement.ChildAttr("#car_memberid", "value")
	if len(memberIdStr) > 0 {
		memberId = cast.ToInt(memberIdStr)
	}
	return memberId
}

// 不确定
func GetIsOutSite(detailPageHtmlElement *colly.HTMLElement) int {

	// 不确定
	var isOutSite int
	isOutSiteStr := detailPageHtmlElement.ChildAttr("#car_isoutsite", "value")
	if len(isOutSiteStr) > 0 {
		isOutSite = cast.ToInt(isOutSiteStr)
	}
	return isOutSite
}

// 是否已售出
func GetIsSold(detailPageHtmlElement *colly.HTMLElement) int {

	var isSold int
	isSoldStr := detailPageHtmlElement.ChildAttr("#car_isselled", "value")
	if len(isSoldStr) > 0 {
		isSold = cast.ToInt(isSoldStr)
	}
	return isSold
}

// 是否是新车
func GetIsNewCar(detailPageHtmlElement *colly.HTMLElement) int {

	var isNewCar int
	isNewCarStr := detailPageHtmlElement.ChildAttr("#isNewCar", "value")
	if len(isNewCarStr) > 0 {
		isNewCar = cast.ToInt(isNewCarStr)
	}
	return isNewCar
}

// 不确定
func GetIsFactory(detailPageHtmlElement *colly.HTMLElement) int {

	//  不确定
	var isFactory int
	isFactoryStr := detailPageHtmlElement.ChildAttr("#isFactory", "value")
	if len(isFactoryStr) > 0 {
		isFactory = cast.ToInt(isFactoryStr)
	}
	return isFactory
}

// 车龄
func GetCarAge(detailPageHtmlElement *colly.HTMLElement) int {

	var age int
	ageStr := detailPageHtmlElement.ChildAttr("#car_age", "value")
	if len(ageStr) > 0 {
		age = cast.ToInt(ageStr)

	}
	return age
}

// 车辆图片
func GetCarPhotos(detailPageHtmlElement *colly.HTMLElement) string {

	var urls []string
	imgList := detailPageHtmlElement.ChildAttrs("#pic_li>a img", "data-original")

	// 下载车辆图片
	for _, val := range imgList {
		newVal := "https:" + val
		urls = append(urls, newVal)
	}
	urlsStr := strings.Join(urls, ",")
	return urlsStr
}
