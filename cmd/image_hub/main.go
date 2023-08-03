package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"image_hub/model"
	"image_hub/spiders"

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

	// Define the directory to traverse
	// 头像社
	// dir = "D:/work/wechat_download_data/html/Dump-0421-11-15-39"

	// 头像有点好看
	// dir = "D:/work/wechat_download_data/html/Dump-0422-20-45-37"

	// 头像即新欢
	// dir = "D:/work/wechat_download_data/html/Dump-0422-20-54-12"

	// 头像库
	dir = "D:/work/wechat_download_data/html/Dump-0423-11-39-39"
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

	params := make(map[string]interface{})
	// 计算目录下的html内的img标签重复的data-src
	dataSrcRepeat, err := spiders.GetImageDataSrcRepeat(dir)
	if err != nil {
		panic(err)
	}
	params["dataSrcRepeat"] = dataSrcRepeat

	// request local files
	// https://github.com/gocolly/colly/blob/master/_examples/local_files/local_files.go
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	// Collector
	// 爬取本地文件时,不用设置AllowedDomains
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.WithTransport(t)

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 3,
		RandomDelay: 16 * time.Second,
	})

	// create a request queue with 2 consumer threads
	q, _ := queue.New(
		20, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 1000000}, // Use default queue storage
	)

	c.OnRequest(func(r *colly.Request) {

		fmt.Printf("=============== c.OnRequest. url: %v\n", r.URL.String())

		urlType := r.Ctx.Get(spiders.UrlTypeKey)

		isEmpty := q.IsEmpty()
		size, err := q.Size()
		if err != nil {
			log.Infof("Queue.Size() return an error: %v", err)
		}
		threads := q.Threads
		log.Infof("OnRequest: Req.ID: %d, urlType:%s, Req.URL: %s, q.IsEmpty: %+v, q.Size: %d, q.Threads: %d\n", r.ID, urlType, r.URL, isEmpty, size, threads)
	})

	c.OnResponse(func(r *colly.Response) {

		fmt.Printf("=============== c.OnResponse. url: %v\n", r.Request.URL.String())
		urlType := r.Ctx.Get(spiders.UrlTypeKey)
		isEmpty := q.IsEmpty()
		size, err := q.Size()
		if err != nil {
			log.Infof("Queue.Size() return an error: %v", err)
		}
		threads := q.Threads
		log.Infof("OnResponse: Req.ID: %d, urlType:%s, Req.URL: %s, Res.Body.len: %d bytes, q.IsEmpty: %+v, q.Size: %d, q.Threads: %d\n", r.Request.ID, urlType, r.Request.URL, len(r.Body), isEmpty, size, threads)
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {

		urlType := e.Response.Ctx.Get(spiders.UrlTypeKey)
		fmt.Printf("=============== c.OnHTML: [%d]%s, %s\n", e.Request.ID, urlType, e.Request.URL)

		onePageSpider := spiders.NewOnePage(spiders.OnePage)
		err := onePageSpider.Process(onePageSpider, q, e, params)
		if err != nil {
			log.Errorf("onePageSpider.Process failed. err: %s\n", err)
			fmt.Printf("onePageSpider.Process failed. err: %s\n", err)
		}
	})

	// OnScraped中获取的urlType参数错误,先忽略
	c.OnScraped(func(r *colly.Response) {

		// urlType := r.Ctx.Get(UrlTypeKey)
		// log.Infof("OnScraped: [%d]%s,%s\n", r.Request.ID, urlType, r.Request.URL)
	})

	c.OnError(func(r *colly.Response, err error) {

		urlType := r.Ctx.Get(spiders.UrlTypeKey)
		fmt.Printf("===============  OnError: [%d]%s, %s, %v\n", r.Request.ID, urlType, r.Request.URL, err)
		log.Infof("OnError: [%d]%s, %s, %v\n", r.Request.ID, urlType, r.Request.URL, err)
	})

	// Determine which spider to use based on the file count and file name
	var spider spiders.Spider
	onePageSpider := spiders.NewOnePage(spiders.OnePage)

	// 遍历目录D:\work\wechat_download_data\html\Dump-0421-11-15-39下的所有html文件
	// html文件名规则为："%Y%m%d_%H%M%S"_"序号.html", 例如: 20230109_111900_1.html

	// Define the regular expression to match the file names
	re := regexp.MustCompile(`(\d{8}_\d{6})_(\d+)\.html`)

	// Traverse the directory and process each file
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Check if the file name matches the regular expression
		matches := re.FindStringSubmatch(d.Name())
		if len(matches) != 3 {
			return nil
		}

		// Extract the date and page number from the file name
		dateStr := matches[1]
		pageNumStr := matches[2]
		fmt.Printf("dateStr: %+v\n", dateStr)
		fmt.Printf("pageNumStr: %+v\n", pageNumStr)
		fmt.Printf("path: %+v\n", path)

		// 打开本地 HTML 文件
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		htmlBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		htmlStr := string(htmlBytes)

		// html内的图片类型(头像，壁纸，背景图，表情包)
		imageTypes, _ := spiders.GetHtmlImageTypes(htmlStr)
		fmt.Printf("imageTypes: %+v\n", imageTypes)

		if len(imageTypes) > 0 {
			spider = onePageSpider
		} else {
			log.Warnf("no matching spider found for file %s", d.Name())
			fmt.Printf("==== no matching spider found for file %s", d.Name())
			return nil
		}

		isAd := spiders.IsAd(htmlStr)
		if isAd {
			log.Warnf("is ad file: %s", d.Name())
			return nil
		}

		// 替换 \ 为 /
		// D:\work\wechat_download_data\html\test\20220526_111900_1.html
		// D:/work/wechat_download_data/html/test/20220526_111900_1.html
		params["path"] = strings.ReplaceAll(path, "\\", "/")

		// Process the file with the selected spider
		fmt.Printf("==============  spider.AddReqToQueue. spider: %+v,  path: %+v\n", spider.GetName(), params["path"])
		err = spider.AddReqToQueue(q, nil, params)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Printf("filepath.Walk err: %+v\n", err)
		panic(err)
	}
	q.Run(c)
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
