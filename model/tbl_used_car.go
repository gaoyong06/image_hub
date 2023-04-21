/*
 * @Author: gaoyong gaoyong06@qq.com
 * @Date: 2023-03-09 21:54:28
 * @LastEditors: gaoyong gaoyong06@qq.com
 * @LastEditTime: 2023-03-19 18:06:58
 * @FilePath: \car_hub\model\car_info.go
 * @Description: 二手车信息
 */
package model

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// 目标页面1
// https://www.che168.com/dealer/462752/43762173.html?pvareaid=100519&userpid=340000&usercid=340100&offertype=&offertag=0&activitycartype=0#pos=1#page=14#rtype=10#isrecom=0#filter=29#module=10#refreshid=0#recomid=0#queryid=1678329929$B$6e930327-cb2d-4d0d-94b4-eb8674d59c3b$74284#cartype=70

// 二手车价格
// https://pinguapi.che168.com/v1/auto/usedcarassess.ashx?_appid=pc.pingu&_appversion=2.08v&pid=340000&cid=340100&mileage=80000&firstregtime=2016/03&specid=23367&mark=uahp10031&encoding=gb2312&infoid=43762173&callback=dtcommon.loadReferencePriceCallBack

// 新车价格
// https://apiassess.che168.com/api/NewCarPriceInTax.ashx?_callback=dtcommon.load4SPriceCallBack&_appid=2sc&pid=0&specid=23367&cid=340100

// 车辆口碑
// https://yccacheapigo.che168.com/api/carinfo/specprc?callback=dtcommon.getKouBeiTabInfo&specid=23367

// 相似车辆
// https://tj.che168.com/sim/cardetailrecom?callback=getRecommendCarlistCallback&_appid=pc&infoid=43762173&id=6e930327-cb2d-4d0d-94b4-eb8674d59c3b

// 免费咨询车况
// https://callcenterapi.che168.com/CallCenterApi/v100/BindingNumber.ashx?_appid=2sc.pc&_callback=callnumber.getxphonenumbercallback&fromtype=0&infoid=43762173&uniqueid=6e930327-cb2d-4d0d-94b4-eb8674d59c3b&ts=1678370082&_sign=4EEA7315C28862F95E4B2EC4303ADE92&sessionid=6e930327-cb2d-4d0d-94b4-eb8674d59c3b&detailpageurl=https%253A%2F%2Fwww.che168.com%2Fdealer%2F462752%2F43762173.html%253Fpvareaid%253D100519%2526userpid%253D340000%2526usercid%253D340100%2526offertype%253D%2526offertag%253D0%2526activitycartype%253D0%2523pos%253D1%2523page%253D14%2523rtype%253D10%2523isrecom%253D0%2523filter%253D29%2523module%253D10%2523refreshid%253D0%2523recomid%253D0%2523queryid%253D1678329929%2524B%25246e930327-cb2d-4d0d-94b4-eb8674d59c3b%252474284%2523cartype%253D70&detailpageref=https%253A%2F%2Fwww.che168.com%2Fhefei%2Fa0_0msdgscncgpi1ltocsp14exx0%2F%253Fpvareaid%253D102179&adfrom=0&queryid=1678329929%24b%246e930327-cb2d-4d0d-94b4-eb8674d59c3b%2474284&dealerid=462752&cartype=70&visit_info=6e930327-cb2d-4d0d-94b4-eb8674d59c3b%7C%7C05DF66D7-F89F-4141-B855-4FE4303D1FD1%7C%7C20200618%7C%7C01%7C%7C104%7C%7C-1&offertype=&activitycartype=0&sourceid=0

// 店铺评论？
// https://dianpingapir.che168.com/public/getshopcommentlist?_appid=2sc.pc&_callback=getshopcommentlist&dealerid=462752&infoid=43762173

// 目标页面2
// https://www.che168.com/dealer/449520/47332851.html?pvareaid=100519&userpid=0&usercid=0&offertype=&offertag=0&activitycartype=0#pos=49#page=3#rtype=10#isrecom=1#filter=29#module=10#refreshid=0#recomid=0#queryid=1678284753$B$6e930327-cb2d-4d0d-94b4-eb8674d59c3b$45828#cartype=30

// 配置亮点#
// https://apipcmusc.che168.com/v1/car/getusedcaroptiondata?callback=getUsedCarOptionDataCallback&_appid=2sc.m&infoid=47332851

// 交易标签？
// https://dianpingapir.che168.com/public/getdealertags?_appid=2sc.pc&_callback=getdealertags&dealerid=449520

// 交易排名？
// https://dianpingapir.che168.com/public/getdealerscore?_appid=2sc.pc&_callback=getdealerscore&dealerid=449520

// 各种指标得分
// https://yccacheapigo.che168.com/api/carinfo/getscorebyspecid?callback=dtcommon.getScoreCallbackKoubei&specid=34614

// 目标页面3
// https://www.che168.com/dealer/332268/47143633.html?pvareaid=100519&userpid=0&usercid=0&offertype=&offertag=0&activitycartype=0#pos=2#page=3#rtype=10#isrecom=1#filter=29#module=10#refreshid=0#recomid=29281468#queryid=1678284753$B$6e930327-cb2d-4d0d-94b4-eb8674d59c3b$45828#cartype=30

// 收藏数量
// https://www.che168.com/Handler/CarDetail_v3/Favorite.ashx?action=isfavorite&id=47143633&cityId=500100&specId=47664&mid=0

// 价格
// 目标页面 https://www.che168.com/dealer/332268/47143633.html?pvareaid=100519&userpid=0&usercid=0&offertype=&offertag=0&activitycartype=0#pos=2#page=3#rtype=10#isrecom=1#filter=29#module=10#refreshid=0#recomid=29281468#queryid=1678284753$B$6e930327-cb2d-4d0d-94b4-eb8674d59c3b$45828#cartype=30
// https://pinguapi.che168.com/v1/auto/usedcarassess.ashx?_appid=pc.pingu&_appversion=2.07v&pid=500000&cid=500100&mileage=6700&firstregtime=2020/12&specid=47664&mark=uahp10032&encoding=gb2312&infoid=47143633&callback=dtcommon.loadReferencePriceCallBack

// VIN码
// 出险记录 WAUG********39475
// https://www.che168.com/insurance/index.aspx?pvareaid=110206&infoid=46333726&vincode=0E57818D5E7754E5D5D1E73E57069367A6FB25B23E31CA46

// VIN码
// 车辆维修记录 WAUGFEF56MA03****
// https://www.che168.com/maintenance/vincodesearch.html?source=11&pvareaid=108871&vincode=0E57818D5E7754E5D5D1E73E57069367A6FB25B23E31CA46&infoid=46333726&dealerid=85869&seriesid=538&specId=48666

// WAUGFEF56MA039475

// 好像和电池有关
// https://apipcmusc.che168.com/v1/insurance/getcarbatteryreportdata?callback=getcarbatteryreportdatacallback&_appid=2sc.m&vincode=LaK6x4DCbom%2BCIjlAoyhxJtgFIdzgyZL

// 推荐车辆
// https://tj.che168.com/sim/cardetailrecom?callback=getRecommendCarlistCallback&_appid=pc&infoid=47046593&id=6e930327-cb2d-4d0d-94b4-eb8674d59c3b

// 热度+推荐理由
// https://yccacheapigo.che168.com/api/carinfo/getheatrank?callback=getHeatRankCallback&_appid=2sc&seriesid=4851&infoid=47046593

// 车源价格同车系价格占比 -> 该车价格低于<u id="pricenum">xx%</u>的同系车辆
// https://www.che168.com/handlercore/CarInfo/GetSeriesCarPrice.ashx?v=20210826&infoid=47046593

// 新车价格+税
// https://cacheapigo.che168.com/quoted/dealerminpricebyspec.ashx?specid=48666&cityid=510100&callback=loadnewprice.load4SPriceCallBack&_=1678841164293

// 信息很多，眼下只关注
// 基本信息
// 图片
// 配置参数

// 二手车信息
// https://www.che168.com/dealer/332268/47143633.html?pvareaid=100519&userpid=0&usercid=0&offertype=&offertag=0&activitycartype=0#pos=2#page=3#rtype=10#isrecom=1#filter=29#module=10#refreshid=0#recomid=29281468#queryid=1678284753$B$6e930327-cb2d-4d0d-94b4-eb8674d59c3b$45828#cartype=30
type TblUsedCar struct {
	CarId                  int             `json:"carId" gorm:"primaryKey"` // 车辆id
	Name                   string          `json:"carName"`                 // 车辆名称
	ProvinceId             int             `json:"carProvinceId"`           // 车辆省份Id
	CityId                 int             `json:"carCityId"`               // 车辆城市Id
	MakeId                 int             `json:"makeId"`                  // 车辆品牌Id 取自 car_brandid
	ModelId                int             `json:"modelId"`                 // 车辆系列Id 取自 car_seriesid
	ModelDetailId          int             `json:"modelDetailId"`           // 车辆型号 取自 car_specid
	DealerId               int             `json:"dealerId"`                // 暂未使用 取自 car_dealerid
	MemberId               int             `json:"MemberId"`                // 暂未使用 取自 car_memberid
	IsOutSite              int             `json:"isOutSite"`               // 暂未使用 取自 car_isoutsite
	Vin                    string          `json:"vin"`                     // 通过出险记录,车辆维修记录
	VinCode                string          `json:"vinCode"`                 // 取自 hidvincode
	VinCodeStr             string          `json:"vinCodeStr"`              // 取自 hidvincodestr
	VinCodeSign            string          `json:"vinCodeSign"`             // 取自 hidvincodesign
	IsSold                 int             `json:"isSold"`                  // 是否已售出
	IsNewCar               int             `json:"isNewCar"`                // 是否是新车
	IsFactory              int             `json:"isFactory"`               // 暂未使用 isFactory
	Age                    int             `json:"age"`                     // 车龄
	Tags                   []TblUsedCarTag `json:"tags" gorm:"-"`           // 车辆标签
	TagIds                 string          `json:"tag_ids"`                 // 车辆标签 以逗号分割的标签id列表
	Mileage                int             `json:"mileage"`                 // 表显里程 单位:公里
	FirstRegDate           string          `json:"firstRegDate"`            // 上牌时间-2021年09月/未上牌
	Engine                 string          `json:"engine"`                  // 发动机
	Transmission           string          `json:"transmission"`            // 变速箱档位
	EngineDisplacement     string          `json:"engineDisplacement"`      // 排量
	CLTCEnduranceMileage   string          `json:"CLTCEnduranceMileage"`    // CLTC纯电续航里程
	NEDCEnduranceMileage   string          `json:"NEDCEnduranceMileage"`    // NEDC纯电续航里程
	WLTCEnduranceMileage   string          `json:"WLTCEnduranceMileage"`    // WLTC纯电续航里程
	StandardFastCharge     string          `json:"standardFastCharge"`      // 标准快充
	StandardSlowCharge     string          `json:"standardSlowCharge"`      // 标准慢充
	StandardCapacity       string          `json:"standardCapacity"`        // 标准容量
	EmissionStandard       string          `json:"emissionStandard"`        // 排放标准
	Level                  string          `json:"level"`                   // 车辆级别
	ExteriorColor          string          `json:"exteriorColor"`           // 车身颜色
	FuelGrade              string          `json:"fuelGrade"`               // 燃油标号
	DriveType              string          `json:"driveType"`               // 驱动方式
	Location               string          `json:"location"`                // 车辆所在地
	AnnualExpireDate       string          `json:"annualExpireDate"`        // 年检到期日期
	InsuranceExpireDate    string          `json:"insuranceExpireYear"`     // 保险到期日期
	WarrantyDate           string          `json:"warrantyDate"`            // 质保到期日期
	TransferOwnershipCount int             `json:"transferOwnershipCount"`  // 过户次数
	IsIncludeTransferFee   string          `json:"isIncludeTransferFee"`    // 是否包括过户费
	UsedCarPrice           int64           `json:"usedCarPrice"`            // 二手车销售价格 单位:元 car_price
	NewCarPrice            int             `json:"newCarPrice"`             // 新车价格 单位:元
	Tax                    int             `json:"tax"`                     // 新车税金 单位:元
	NewCarPriceIncludeTax  int             `json:"newCarPriceIncludeTax"`   // 新车含税价 单位:元  从https://apiassess.che168.com/api/NewCarPriceInTax.ashx 取到 newcarprice
	UsedCarPriceRefMin     int             `json:"usedCarPriceRefMin"`      // 二手车参考价-最低 单位:元
	UsedCarPriceRefMax     int             `json:"usedCarPriceRefMax"`      // 二手车参考价-最高 单位:元
	SellerPhoneNumber      string          `json:"sllerPhoneNumber"`        // 卖家电话
	FuelType               string          `json:"fuelType"`                // 燃料类型 纯电动
	Photos                 string          `json:"photo" gorm:"type:text"`  // 车辆图片 以图片url逗号分割字符串
	FavoriteCount          int             `json:"favoriteCount"`           // 收藏数量
	OriginUrl              string          `json:"originUrl"`               // 数据来源url 车辆详情页Url
	PublicDate             string          `json:"publicDate"`              // 发布日期
	CreatedAt              time.Time       `json:"createdAt"`               // 创建时间
	UpdatedAt              time.Time       `json:"updatedAt"`               // 最后修改时间
	DeletedAt              time.Time       `json:"deletedAt"`               // 删除时间
}

func GetTblCarModel() *TblUsedCar {
	return &TblUsedCar{}
}

func (t *TblUsedCar) TableName() string {
	return "tbl_used_car"
}

func (t *TblUsedCar) FirstOrCreate() (int, error) {

	// https://gorm.io/zh_CN/docs/advanced_query.html
	condModel := TblUsedCar{CarId: t.CarId}
	assignModel := TblUsedCar{
		CarId:                  t.CarId,
		Name:                   t.Name,
		ProvinceId:             t.ProvinceId,
		CityId:                 t.CityId,
		MakeId:                 t.MakeId,
		ModelId:                t.ModelId,
		ModelDetailId:          t.ModelDetailId,
		DealerId:               t.DealerId,
		MemberId:               t.MemberId,
		IsOutSite:              t.IsOutSite,
		Vin:                    t.Vin,
		VinCode:                t.VinCode,
		VinCodeStr:             t.VinCodeStr,
		VinCodeSign:            t.VinCodeSign,
		IsSold:                 t.IsSold,
		IsNewCar:               t.IsNewCar,
		IsFactory:              t.IsFactory,
		Age:                    t.Age,
		TagIds:                 t.TagIds,
		Mileage:                t.Mileage,
		FirstRegDate:           t.FirstRegDate,
		Engine:                 t.Engine,
		Transmission:           t.Transmission,
		EngineDisplacement:     t.EngineDisplacement,
		CLTCEnduranceMileage:   t.CLTCEnduranceMileage,
		NEDCEnduranceMileage:   t.NEDCEnduranceMileage,
		WLTCEnduranceMileage:   t.WLTCEnduranceMileage,
		StandardFastCharge:     t.StandardFastCharge,
		StandardSlowCharge:     t.StandardSlowCharge,
		StandardCapacity:       t.StandardCapacity,
		EmissionStandard:       t.EmissionStandard,
		Level:                  t.Level,
		ExteriorColor:          t.ExteriorColor,
		FuelGrade:              t.FuelGrade,
		DriveType:              t.DriveType,
		Location:               t.Location,
		AnnualExpireDate:       t.AnnualExpireDate,
		InsuranceExpireDate:    t.InsuranceExpireDate,
		WarrantyDate:           t.WarrantyDate,
		TransferOwnershipCount: t.TransferOwnershipCount,
		IsIncludeTransferFee:   t.IsIncludeTransferFee,
		UsedCarPrice:           t.UsedCarPrice,
		NewCarPrice:            t.NewCarPrice,
		Tax:                    t.Tax,
		NewCarPriceIncludeTax:  t.NewCarPriceIncludeTax,
		UsedCarPriceRefMin:     t.UsedCarPriceRefMin,
		UsedCarPriceRefMax:     t.UsedCarPriceRefMax,
		SellerPhoneNumber:      t.SellerPhoneNumber,
		FuelType:               t.FuelType,
		Photos:                 t.Photos,
		FavoriteCount:          t.FavoriteCount,
		OriginUrl:              t.OriginUrl,
		PublicDate:             t.PublicDate,
	}

	err := DB.Table(t.TableName()).Where(condModel).Assign(assignModel).FirstOrCreate(t).Error
	if err != nil {
		log.Errorf("TblUsedCar CreateOrUpdate failed. err: %+v\n", err)
	}

	return t.CarId, err
}

func (t *TblUsedCar) Update(carId int, updates map[string]interface{}) error {

	if len(updates) < 1 {
		return nil
	}

	err := DB.Table(t.TableName()).Where("car_id=?", carId).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TblUsedCar) UpdateByVinCodeSign(vinCodeSign string, updates map[string]interface{}) error {

	if len(updates) < 1 {
		return nil
	}

	err := DB.Table(t.TableName()).Where("vin_code_sign=?", vinCodeSign).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TblUsedCar) UpdateByCityIdAndModelDetailId(cityId, modelDetailId int, updates map[string]interface{}) error {

	if len(updates) < 1 {
		return nil
	}

	err := DB.Table(t.TableName()).Where("city_id=? AND model_detail_id=?", cityId, modelDetailId).Updates(updates).Error
	if err != nil {
		return err
	}
	return nil
}
