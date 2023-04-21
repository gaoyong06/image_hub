/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-03-09 11:54:15
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-03-24 21:09:10
 * @FilePath: \car_hub\spiders\detail.go
 * @Description: 二手车详情页数据抓取
 */
package spiders

import (
	"car_hub/model"
	"car_hub/pkg/utils"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

type DetailPageSpider struct {
	Name string
}

func (s *DetailPageSpider) GetName() string {
	return s.Name
}

func (s *DetailPageSpider) SetName(name string) {
	s.Name = name
}

// 抓取详情页
func (s *DetailPageSpider) AddReqToQueue(q *queue.Queue, listPageHtmlElement *colly.HTMLElement, baseUrl string) error {

	// 获取列表页的所有详情页url
	detailPageHrefs := listPageHtmlElement.ChildAttrs(".carinfo", "href")
	var detailPageUrl string

	for _, detailPageHref := range detailPageHrefs {

		log.Infof("DetailPageSpider AddReqToQueue range. detailPageHref: %+v, url: %+v\n", detailPageHref, listPageHtmlElement.Request.URL.String())
		if strings.HasPrefix(detailPageHref, "//www.che168.com") {

			detailPageUrl = "https:" + detailPageHref
		} else if strings.HasPrefix(detailPageHref, "/dealer") {

			detailPageUrl = "https://www.che168.com" + detailPageHref
		} else {

			log.Infof("DetailPageSpider AddReqToQueue detailPageHref failed. detailPageHref: %+v, url: %+v\n", detailPageHref, listPageHtmlElement.Request.URL.String())
		}

		// 继续抓取详情页
		if _, ok := visited.Get(detailPageUrl); !ok {

			visited.Set(detailPageUrl, true)

			detailPageReq, err := listPageHtmlElement.Request.New("GET", detailPageUrl, nil)
			if err != nil {

				log.Errorf("DetailPageSpider AddReqToQueue Request.New failed. error: %+v\n", err)
				return err
			}

			detailPageReq.Ctx.Put(UrlTypeKey, UrlTypeUsedCarDetailPage)
			q.AddRequest(detailPageReq)
			log.Infof("DetailPageSpider AddReqToQueue success. detailPageUrl: %+v\n", detailPageUrl)
		}
	}

	return nil
}

// 解析详情页页面结构,获取二手车车辆信息
func (s *DetailPageSpider) ParseData(q *queue.Queue, e *colly.HTMLElement, baseUrl string) (*model.TblUsedCar, error) {

	var err error
	model := model.GetTblCarModel()

	// 解析详情页二手车信息

	// 二手车id
	model.CarId = GetCarId(e)

	// 二手车名称
	model.Name = GetName(e)

	// 二手车价格
	model.UsedCarPrice = GetPrice(e)

	// 是否包含过户费
	model.IsIncludeTransferFee = GetIsIncludeTransferFee(e)

	// 二手车所属省id
	model.ProvinceId = GetProvinceIdFromDetailPage(e)

	// 二手车所属城市id
	model.CityId = GetCityId(e)

	// 二手车品牌id
	model.MakeId = GetMakeId(e)

	// 二手车车系id
	model.ModelId = GetModelId(e)

	// 二手车车型id
	model.ModelDetailId = GetModelDetailIdFromDetailPage(e)

	// 二手车卖家id
	model.DealerId = GetDealerId(e)

	// 不确定
	model.MemberId = GetMemberId(e)

	// 不确定
	model.IsOutSite = GetIsOutSite(e)

	// 二手车VIN号相关, VinCode, VinCodeStr, VinCodeSign 都是编码后的VIN
	model.VinCode = e.ChildAttr("#hidvincode", "value")
	model.VinCodeStr = e.ChildAttr("#hidvincodestr", "value")

	// 签名后的VIN码
	model.VinCodeSign = GetVinCodeSign(e)

	// 是否已售出
	model.IsSold = GetIsSold(e)

	// 是否是新车
	model.IsNewCar = GetIsNewCar(e)

	//  不确定
	model.IsFactory = GetIsFactory(e)

	//  车龄
	model.Age = GetCarAge(e)

	// 标签
	model.Tags = GetCarTags(e)

	// 定位到ul元素
	selUl := e.DOM.Find(".all-basic-content.fn-clear").Find(".basic-item-ul")

	// 上牌时间
	model.FirstRegDate = GetFirstRegDate(e)

	// 表显里程
	model.Mileage = getMileage(e)

	// 变速箱
	selLi := selUl.Eq(0).Find("li").Eq(2)
	text := utils.GetNodeTextOne(selLi)
	title := utils.RemoveSpace(selLi.Find(".item-name").Text())
	if title == "变速箱" {
		model.Transmission = text
	}

	// 排放标准
	selLi = selUl.Eq(0).Find("li").Eq(3)
	text = utils.GetNodeTextOne(selLi)
	title = utils.RemoveSpace(selLi.Find(".item-name").Text())
	if title == "排放标准" {
		model.EmissionStandard = text
	}
	// 新能源-燃料类型
	if title == "燃料类型" {
		model.FuelType = text
	}

	// 排量
	selLi = selUl.Eq(0).Find("li").Eq(4)
	text = utils.GetNodeTextOne(selLi)
	title = utils.RemoveSpace(selLi.Find(".item-name").Text())
	if title == "排量" {
		model.EngineDisplacement = text
	}
	// 新能源-xxxx纯电续航里程
	if title == "CLTC纯电续航里程" {
		model.CLTCEnduranceMileage = text
	}
	if title == "NEDC纯电续航里程" {
		model.NEDCEnduranceMileage = text
	}
	if title == "WLTC纯电续航里程" {
		model.WLTCEnduranceMileage = text
	}

	// 发布时间
	selLi = selUl.Eq(0).Find("li").Eq(5)
	text = utils.GetNodeTextOne(selLi)
	title = utils.RemoveSpace(selLi.Find(".item-name").Text())
	if title == "发布时间" {
		model.PublicDate = text
	}

	// 年检到期
	selLi = selUl.Eq(1).Find("li").Eq(0)
	text = utils.GetNodeTextOne(selLi)
	title = utils.RemoveSpace(selLi.Find(".item-name").Text())
	if title == "年检到期" {
		model.AnnualExpireDate = text
	}

	// 保险到期
	selLi = selUl.Eq(1).Find("li").Eq(1)
	text = utils.GetNodeTextOne(selLi)
	title = utils.RemoveSpace(selLi.Find(".item-name").Text())
	if title == "保险到期" {
		model.InsuranceExpireDate = text
	}

	// 质保到期
	selLi = selUl.Eq(1).Find("li").Eq(2)
	text = utils.GetNodeTextOne(selLi)
	title = utils.RemoveSpace(selLi.Find(".item-name").Text())
	if title == "质保到期" {
		model.WarrantyDate = text
	}

	// 过户次数
	model.TransferOwnershipCount = GetTransferOwnershipCount(e)

	// 所在地
	selLi = selUl.Eq(1).Find("li").Eq(5)
	text = utils.GetNodeTextOne(selLi)
	title = utils.RemoveSpace(selLi.Find(".item-name").Text())
	if title == "所在地" {
		model.Location = text
	}

	// 新能源-标准快充
	selLi = selUl.Eq(1).Find("li").Eq(6)
	if selLi != nil {
		text = utils.GetNodeTextOne(selLi)
		title = utils.RemoveSpace(selLi.Find(".item-name").Text())
		if title == "标准快充" {
			model.StandardFastCharge = text
		}
	}

	// 发动机
	selLi = selUl.Eq(2).Find("li").Eq(0)
	text = utils.GetNodeTextOne(selLi)
	title = utils.RemoveSpace(selLi.Find(".item-name").Text())

	if title == "发动机" {
		model.Engine = text
	}

	// 车辆级别
	selLi = selUl.Eq(2).Find("li").Eq(1)
	text = utils.GetNodeTextOne(selLi)
	title = selLi.Find(".item-name").Text()
	if title == "车辆级别" {
		model.Level = text
	}

	// 车身颜色
	selLi = selUl.Eq(2).Find("li").Eq(2)
	text = utils.GetNodeTextOne(selLi)
	title = utils.RemoveSpace(selLi.Find(".item-name").Text())
	if title == "车身颜色" {
		model.ExteriorColor = text
	}

	// 燃油标号,新能源-驱动方式
	selLi = selUl.Eq(2).Find("li").Eq(3)
	text = utils.GetNodeTextOne(selLi)
	title = utils.RemoveSpace(selLi.Find(".item-name").Text())
	if title == "燃油标号" {
		model.FuelGrade = text
	}
	// 新能源-驱动方式
	if title == "驱动方式" {
		model.DriveType = text
	}

	// 燃油车-驱动方式
	selLi = selUl.Eq(2).Find("li").Eq(4)
	text = utils.GetNodeTextOne(selLi)
	title = utils.RemoveSpace(selLi.Find(".item-name").Text())
	if title == "驱动方式" {
		model.DriveType = text
	}

	// 新能源-标注容量
	selLi = selUl.Eq(2).Find("li").Eq(5)
	text = utils.GetNodeTextOne(selLi)
	title = utils.RemoveSpace(selLi.Find(".item-name").Text())
	if title == "标准容量" {
		model.StandardCapacity = text
	}

	// 新能源-标准慢充
	selLi = selUl.Eq(2).Find("li").Eq(6)
	text = utils.GetNodeTextOne(selLi)
	title = utils.RemoveSpace(selLi.Find(".item-name").Text())
	if title == "标准慢充" {
		model.StandardSlowCharge = text
	}

	// 是否包括过户费
	text = e.ChildText(".price-transfer.ndy")
	model.IsIncludeTransferFee = text

	// // TODO: 收藏数量
	// // model.FavoriteCount

	// 数据来源url 车辆详情页Url
	model.OriginUrl, err = url.QueryUnescape(e.Request.URL.String())

	// 车辆图片
	model.Photos = GetCarPhotos(e)
	return model, err
}

// 业务处理(如：保存)及后续的行为(如：继续发起下一个层级的请求)
func (s *DetailPageSpider) Process(q *queue.Queue, e *colly.HTMLElement, baseUrl string) error {

	// 解析详情页页面结构
	tblModel, err := s.ParseData(q, e, BaseUrl)
	if err != nil {
		log.Errorf("DetailPageSpider ParseData failed. error: %+v, url: %+v\n", err, e.Request.URL.String())
		return err
	}

	// spew.Dump(tblModel)

	// 保存标签
	usedCarTags := tblModel.Tags
	if len(usedCarTags) > 0 {

		// 添加至tbl_used_tag
		var tagIds []string
		for _, tag := range usedCarTags {

			tagId, err := tag.FirstOrCreate()
			if err != nil {

				panic(err)
			}
			tagIds = append(tagIds, cast.ToString(tagId))
		}
		tblModel.TagIds = strings.Join(tagIds, ",")
	}

	carId, err := tblModel.FirstOrCreate()
	if err != nil {
		log.Errorf("DetailPageSpider UsedCarCreate failed. carId: %d, error: %s \n", carId, err.Error())
		return err
	} else {
		log.Infof("UsedCarCreate success. carId:%d \n", carId)
	}

	// 增加下载车辆图片请求至队列-done
	// 暂时注释掉
	// err = addDownloadImgReqToQueue(q, e, tblModel)
	// if err != nil {
	// 	log.Errorf("DetailPageSpider addDownloadImgReqToQueue failed. carId: %d, error: %s \n", carId, err.Error())
	// 	return err
	// }

	// 抓取二手车价格 - done
	// usedCarPriceSpider := UsedCarPriceSpider{Name: UrlTypeUsedCarPriceApi}
	// err = usedCarPriceSpider.AddReqToQueue(q, e, BaseUrl)
	// if err != nil {
	// 	log.Errorf("DetailPageSpider usedCarPriceSpider.AddReqToQueue failed. err: %+v\n", err)
	// 	return err
	// }

	// 抓取4s店新车含税价格 - done
	// newCarPriceIncludeTaxSpider := NewCarPriceIncludeTaxSpider{Name: UrlTypeNewCarPriceIncludeTaxApi}
	// err = newCarPriceIncludeTaxSpider.AddReqToQueue(q, e, BaseUrl)
	// if err != nil {
	// 	log.Errorf("DetailPageSpider newCarPriceIncludeTaxSpider.AddReqToQueue failed. err: %+v\n", err)
	// 	return err
	// }

	// 抓取二手车热度指数，推荐理由 - done
	// usedCarRankSpider := UsedCarRankSpider{Name: UrlTypeUsedCarRankApi}
	// err = usedCarRankSpider.AddReqToQueue(q, e, BaseUrl)
	// if err != nil {
	// 	log.Errorf("DetailPageSpider usedCarRankSpider.AddReqToQueue failed. err: %+v\n", err)
	// 	return err
	// }

	// 抓取二手车卖家电话号码 - done
	// sellerPhoneSpider := SellerPhoneSpider{Name: UrlTypeSellerPhoneApi}
	// err = sellerPhoneSpider.AddReqToQueue(q, e, BaseUrl)
	// if err != nil {
	// 	log.Errorf("DetailPageSpider sellerPhoneSpider.AddReqToQueue failed. err: %+v\n", err)
	// 	return err
	// }

	// // 抓取车辆维修保养记录查询获取VIN码 - done
	// vinCodeSearchSpider := VinCodeSearchSpider{Name: UrlTypeSellerPhoneApi}
	// err = vinCodeSearchSpider.AddReqToQueue(q, e, BaseUrl)
	// if err != nil {
	// 	log.Errorf("DetailPageSpider vinCodeSearchSpider.AddReqToQueue failed. err: %+v\n", err)
	// 	return err
	// }

	// 抓取配置亮点API
	// UsedCarOptionSpider := UsedCarOptionSpider{Name: UrlTypeUsedCarOptionApi}
	// err = UsedCarOptionSpider.AddReqToQueue(q, e, BaseUrl)
	// if err != nil {
	// 	log.Errorf("DetailPageSpider UsedCarOptionSpider.AddReqToQueue failed. err: %+v\n", err)
	// 	return err
	// }

	// 抓取配置参数页
	// configPageSpider := ConfigPageSpider{Name: UrlTypeUsedCarOptionApi}
	// err = configPageSpider.AddReqToQueue(q, e, BaseUrl)
	// if err != nil {
	// 	log.Errorf("DetailPageSpider configPageSpider.AddReqToQueue failed. err: %+v\n", err)
	// 	return err
	// }

	return nil
}

// 增加下载车辆图片请求至队列
func addDownloadImgReqToQueue(q *queue.Queue, e *colly.HTMLElement, m *model.TblUsedCar) error {

	photos := strings.Split(m.Photos, ",")

	// 下载车辆图片
	for _, photo := range photos {

		req, err := e.Request.New("GET", photo, nil)
		if err != nil {
			log.Errorf("addDownloadImgReqToQueue failed. error: %s \n", err.Error())
			return err
		}
		q.AddRequest(req)
	}

	// isEmpty := q.IsEmpty()
	// size, err := q.Size()
	// if err != nil {
	// 	log.Infof("Queue.Size() return an error: %v", err)
	// }
	// threads := q.Threads
	// log.Infof("addDownloadImgReqToQueue. q.IsEmpty: %+v, q.Size: %d, q.Threads: %d ==================\n", isEmpty, size, threads)
	return nil
}
