package model

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	// 数据库连接
	DB           *gorm.DB
	contentNerDB *gorm.DB

	HostDefault     string = "127.0.0.1"
	PortDefault     string = "3306"
	UserDefault     string = "root"
	PasswdDefault   string = "root"
	DatabaseDefault string = "image_hub"
)

func Init() {

	var err error
	DB, err = openDB("default")
	if err != nil {
		panic(fmt.Sprintf("openDB default failed. error: %s\n", err))
	}

	contentNerDB, err = openDB("content_ner_db")
	if err != nil {
		panic(fmt.Sprintf("openDB default failed. error: %s\n", err))
	}
}

func openDB(dbConfig string) (*gorm.DB, error) {

	var err error
	logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			// LogLevel:                  logger.Info, // 日志级别
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false, // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false, // 禁用彩色打印
		},
	)

	namingStrategy := schema.NamingStrategy{
		// TablePrefix:   "tbl_", // table name prefix, table for `User` would be `t_users`
		SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
	}

	// dns := username:password@tcp(host:port)/dbname?charset=utf8&parseTime=True&loc=Local
	dns := getDns(dbConfig)
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		Logger:         logger,
		NamingStrategy: namingStrategy,
	})
	if err != nil {
		panic(fmt.Sprintf("model init gorm.Open failed. error: %s\n", err))
	}

	// https://gorm.cn/zh_CN/docs/generic_interface.html
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("model init db.DB failed. error: %s\n", err))
	}

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func getDns(name string) string {

	host := viper.GetString("mysql." + name + ".host")
	port := viper.GetString("mysql." + name + ".port")
	user := viper.GetString("mysql." + name + ".user")
	passwd := viper.GetString("mysql." + name + ".passwd")
	database := viper.GetString("mysql." + name + ".database")

	if host == "" {
		host = HostDefault
	}
	if port == "" {
		port = PortDefault
	}
	if user == "" {
		user = UserDefault
	}
	if passwd == "" {
		passwd = PasswdDefault
	}
	if database == "" {
		database = DatabaseDefault
	}

	dsn := user + ":" +
		passwd + "@" +
		"tcp(" + host + ":" +
		port + ")/" +
		database +
		"?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local"

	return dsn
}
