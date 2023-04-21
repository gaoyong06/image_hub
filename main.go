package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"car_hub/fake"
	"car_hub/model"
	"car_hub/pkg/utils"
	"car_hub/spiders"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	MODE    string
	VERSION string

	serviceName       string
	configFile        string
	defaultConfigFile = "conf/config.yaml"
)

func main() {

	if MODE == "" {
		MODE = "debug"
	}

	err := Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	Run()
}

// 初始化
func Init() error {

	err := initFlag()
	if err != nil {
		return err
	}

	//init config
	err = initConfig(configFile)
	if err != nil {
		return errors.New("config init failed, err:" + err.Error())
	}

	serviceName = viper.GetString("server.server_name")
	err = initLog()
	if err != nil {
		return err
	}

	model.Init()

	return nil
}

func Run() {

	// 获取可被抓取的域名
	domains := strings.Split(spiders.Domains, ",")

	// 使用go_proxy_pool 得到的代理ip和port
	// proxyUrls := []string{
	// 	"http://117.160.250.133:9999",
	// 	"http://117.159.15.99:9091",
	// 	"http://116.113.68.130:9091",
	// 	"http://221.6.215.202:9091",
	// 	"http://117.160.250.138:8080",
	// 	"http://117.160.250.134:8081",
	// 	"http://117.160.250.131:8899",
	// 	"http://27.15.232.197:9091",
	// 	"http://222.139.221.185:9091",
	// 	"http://42.228.61.245:9091",
	// }

	// rp, err := proxy.RoundRobinProxySwitcher(proxyUrls...)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Collector
	collector := colly.NewCollector(
		colly.AllowedDomains(domains...),
		colly.AllowURLRevisit(),
	)
	collector.SetRequestTimeout(120 * time.Second)
	// collector.SetProxyFunc(rp)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 3,
		RandomDelay: 16 * time.Second,
	})

	// create a request queue with 2 consumer threads
	q, _ := queue.New(
		20, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 1000000}, // Use default queue storage
	)

	// https://www.che168.com/china/list/ 爬虫
	chinaListPageSpider := &spiders.ChinaListPageSpider{Name: spiders.UrlTypeChinaListPage}

	collector.OnRequest(func(r *colly.Request) {

		fake.SetChe168Headers(r)
		urlType := r.Ctx.Get(spiders.UrlTypeKey)

		isEmpty := q.IsEmpty()
		size, err := q.Size()
		if err != nil {
			log.Infof("Queue.Size() return an error: %v", err)
		}
		threads := q.Threads
		log.Infof("OnRequest: Req.ID: %d, urlType:%s, Req.URL: %s, q.IsEmpty: %+v, q.Size: %d, q.Threads: %d\n", r.ID, urlType, r.URL, isEmpty, size, threads)
	})

	collector.OnResponse(func(r *colly.Response) {

		urlType := r.Ctx.Get(spiders.UrlTypeKey)
		isEmpty := q.IsEmpty()
		size, err := q.Size()
		if err != nil {
			log.Infof("Queue.Size() return an error: %v", err)
		}
		threads := q.Threads

		log.Infof("OnResponse: Req.ID: %d, urlType:%s, Req.URL: %s, Res.Body.len: %d bytes, q.IsEmpty: %+v, q.Size: %d, q.Threads: %d\n", r.Request.ID, urlType, r.Request.URL, len(r.Body), isEmpty, size, threads)

		// 限流
		if strings.Contains(r.Request.URL.String(), "www.che168.com/cheerror.html") {
			panic("Oh yeah. limited!")
		}

		// 图片保存
		err = utils.SaveImage(r, spiders.ImageDir)
		if err != nil {
			log.Errorf("utils.SaveImage failed. err: %s\n", err)
		}

		switch urlType {

		// 新车价格API
		case spiders.UrlTypeNewCarPriceApi:
			newCarPriceSpider := spiders.NewCarPriceSpider{Name: spiders.UrlTypeNewCarPriceApi}
			err := newCarPriceSpider.Process(q, r, spiders.BaseUrl)
			if err != nil {
				log.Errorf("newCarPriceSpider.Process failed. err: %s\n", err)
			}

		// 二手车价格API结果返回
		case spiders.UrlTypeUsedCarPriceApi:
			usedCarPriceSpider := &spiders.UsedCarPriceSpider{Name: spiders.UrlTypeUsedCarPriceApi}
			err := usedCarPriceSpider.Process(q, r, spiders.BaseUrl)
			if err != nil {
				log.Errorf("usedCarPriceSpider.Process failed. err: %s\n", err)
			}

		// 4S店新车含税价格API结果返回
		case spiders.UrlTypeNewCarPriceIncludeTaxApi:
			newCarPriceIncludeTaxSpider := &spiders.NewCarPriceIncludeTaxSpider{Name: spiders.UrlTypeNewCarPriceIncludeTaxApi}
			err := newCarPriceIncludeTaxSpider.Process(q, r, spiders.BaseUrl)
			if err != nil {
				log.Errorf("newCarPriceIncludeTaxSpider.Process failed. err: %s\n", err)
			}

		// 二手车在该省最低,最高成交价API结果返回
		case spiders.UrlTypeUsedCarProvincePriceApi:
			newCarPriceIncludeTaxSpider := &spiders.NewCarPriceIncludeTaxSpider{Name: spiders.UrlTypeNewCarPriceIncludeTaxApi}
			err := newCarPriceIncludeTaxSpider.Process(q, r, spiders.BaseUrl)
			if err != nil {
				log.Errorf("newCarPriceIncludeTaxSpider.Process failed. err: %s\n", err)
			}

		// 车辆参数API结果返回
		case spiders.UrlTypeCarParamApi:
			carParamSpider := spiders.CarParamSpider{Name: spiders.UrlTypeCarParamApi}
			err := carParamSpider.Process(q, r, spiders.BaseUrl)
			if err != nil {
				log.Errorf("carParamSpider.Process failed. err: %s\n", err)
			}

		// 车辆配置API结果返回
		case spiders.UrlTypeCarConfigApi:
			carConfigSpider := spiders.CarConfigSpider{Name: spiders.UrlTypeCarConfigApi}
			err := carConfigSpider.Process(q, r, spiders.BaseUrl)
			if err != nil {
				log.Errorf("carConfigSpider.Process failed. err: %s\n", err)
			}

		// 二手车热度指数，推荐理由API结果返回
		case spiders.UrlTypeUsedCarRankApi:
			usedCarRankSpider := &spiders.UsedCarRankSpider{Name: spiders.UrlTypeUsedCarRankApi}
			err := usedCarRankSpider.Process(q, r, spiders.BaseUrl)
			if err != nil {
				log.Errorf("usedCarRankSpider.Process failed. err: %s\n", err)
			}

		// 二手车卖家电话API结果返回
		case spiders.UrlTypeSellerPhoneApi:
			sellerPhoneSpider := &spiders.SellerPhoneSpider{Name: spiders.UrlTypeSellerPhoneApi}
			err := sellerPhoneSpider.Process(q, r, spiders.BaseUrl)
			if err != nil {
				log.Errorf("sellerPhoneSpider.Process failed. err: %s\n", err)
			}

		// 二手车配置亮点API结果返回
		case spiders.UrlTypeUsedCarOptionApi:
			usedCarOptionSpider := &spiders.UsedCarOptionSpider{Name: spiders.UrlTypeUsedCarOptionApi}
			err := usedCarOptionSpider.Process(q, r, spiders.BaseUrl)
			if err != nil {
				log.Errorf("usedCarOptionSpider.Process failed. err: %s\n", err)
			}
		}
	})

	collector.OnHTML("html", func(e *colly.HTMLElement) {

		urlType := e.Response.Ctx.Get(spiders.UrlTypeKey)
		log.Infof("OnHTML: [%d]%s, %s\n", e.Request.ID, urlType, e.Request.URL)

		switch urlType {

		// "全国"列表页url
		case spiders.UrlTypeChinaListPage:

			err := chinaListPageSpider.Process(q, e, spiders.BaseUrl)
			if err != nil {
				log.Errorf("chinaListPageSpider.Process failed. err: %s\n", err)
			}

		// 当前选中的城市列表页url
		case spiders.UrlTypeCityListPage:
			cityListPageSpider := &spiders.CityListPageSpider{
				Name: spiders.UrlTypeChinaListPage,
			}
			err := cityListPageSpider.Process(q, e, spiders.BaseUrl)
			if err != nil {
				log.Errorf("cityListPageSpider.Process failed. err: %s\n", err)
			}

		// 详情页处理
		case spiders.UrlTypeUsedCarDetailPage:
			detailPageSpider := &spiders.DetailPageSpider{
				Name: spiders.UrlTypeUsedCarDetailPage,
			}
			err := detailPageSpider.Process(q, e, spiders.BaseUrl)
			if err != nil {
				log.Errorf("detailPageSpider.Process failed. err: %s\n", err)
			}

		// 判断当前返回html为配置参数页
		case spiders.UrlTypeUsedCarConfigPage:
			configPageSpider := &spiders.ConfigPageSpider{
				Name: spiders.UrlTypeUsedCarConfigPage,
			}
			err := configPageSpider.Process(q, e, spiders.BaseUrl)
			if err != nil {
				log.Errorf("configPageSpider.Process failed. err: %s\n", err)
			}

		// 车辆维修保养记录查询获取VIN码结果返回
		case spiders.UrlTypeUsedCarVinCodeSearchApi:
			vinCodeSearchSpider := &spiders.VinCodeSearchSpider{Name: spiders.UrlTypeUsedCarPriceApi}
			err := vinCodeSearchSpider.Process(q, e, spiders.BaseUrl)
			if err != nil {
				log.Errorf("vinCodeSearchSpider.Process failed. err: %s\n", err)
			}
		}
	})

	// OnScraped中获取的urlType参数错误,先忽略
	// collector.OnScraped(func(r *colly.Response) {

	// 	urlType := r.Ctx.Get(UrlTypeKey)
	// 	log.Infof("OnScraped: [%d]%s,%s\n", r.Request.ID, urlType, r.Request.URL)
	// })

	collector.OnError(func(r *colly.Response, err error) {

		urlType := r.Ctx.Get(spiders.UrlTypeKey)
		log.Infof("OnError: [%d]%s, %s, %v\n", r.Request.ID, urlType, r.Request.URL, err)
	})

	// 开始抓取 https://www.che168.com/china/list/
	err := chinaListPageSpider.AddReqToQueue(q, nil, spiders.BaseUrl)
	if err != nil {
		log.Errorf("main chinaListPageSpider.AddReqToQueue failed. error: %+v\n", err.Error())
		panic(err)
	}
	q.Run(collector)
}

// init flag
func initFlag() error {
	//init params
	h := flag.Bool("h", false, "application help")

	c := flag.String("c", "", "config file, ex: /data/config.yaml")

	flag.Parse()

	if *h {
		flag.PrintDefaults()
		os.Exit(0)
		return nil
	}

	if *c == "" {
		//use default config file
		path, _ := os.Executable()
		binPath := filepath.Dir(path)
		*c = binPath + "/../" + defaultConfigFile
	}

	configFile = *c

	return nil
}

func initConfig(configFile string) error {
	path, _ := os.Executable()
	RootPath := filepath.Dir(path)

	viper.Set("path.root", RootPath)

	configPath := filepath.Dir(configFile)
	fileName := filepath.Base(configFile)

	//设置读取的配置文件
	viper.SetConfigName(fileName)
	//添加读取的配置文件路径
	viper.AddConfigPath(configPath)
	//设置配置文件类型
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func initLog() error {

	logPathPrefix := fmt.Sprintf("%s/%s", viper.GetString("path.root"), viper.GetString("log.prefix"))
	logName := viper.GetString("log.name")
	logLevel := viper.GetInt("log.level")
	maxAge := viper.GetInt("log.maxAge")
	path := fmt.Sprintf("%s%s", logPathPrefix, logName)

	/* 日志轮转相关函数
	`WithLinkName` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	`WithMaxAge 和 WithRotationCount二者只能设置一个
	`WithMaxAge` 设置文件清理前的最长保存时间
	`WithRotationCount` 设置文件清理前最多保存的个数
	 WithMaxAge WithRotationCount 只能存在一个
	*/
	// 下面配置日志每隔 1 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。
	writer, err := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithRotationTime(time.Duration(86400)*time.Second),
		rotatelogs.WithMaxAge(time.Duration(maxAge*24)*time.Hour),
	)
	if err != nil {
		return errors.New("initLog error:" + err.Error())
	}
	log.SetReportCaller(true)
	log.SetFormatter(&nested.Formatter{
		HideKeys:         true,
		NoUppercaseLevel: true,
		ShowFullLevel:    true,
		TimestampFormat:  "2006/01/02 15:04:05",
		NoFieldsSpace:    true,
		NoFieldsColors:   true,
		NoColors:         true,
		CallerFirst:      true,
		TrimMessages:     true,
		FieldsOrder:      []string{"component", "category"},
	})
	log.SetLevel(log.Level(logLevel))
	log.SetOutput(writer)
	return nil
}
