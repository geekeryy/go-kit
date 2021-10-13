package xmysql

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/comeonjy/go-kit/pkg/xlog"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config Mysql配置结构体
type Config struct {
	User        string `json:"user"`          // 用户名
	Password    string `json:"password"`      // 密码
	Host        string `json:"host" `         // 主机地址
	Port        int    `json:"port" `         // 端口号
	Dbname      string `json:"dbname"`        // 数据库名
	MaxIdleConn int    `json:"max_idle_conn"` // 最大空闲连接
	MaxOpenConn int    `json:"max_open_conn"` // 最大活跃连接
	LogLevel    int    `json:"log_level"`     // 日志等级 logger.LogLevel
	Colorful    bool   `json:"colorful"`      // 是否开启彩色日志
}

// New 初始化数据库
func New(config string, xlogger *xlog.Logger) *gorm.DB {
	mysqlConfig := Config{}
	if err := json.Unmarshal([]byte(config), &mysqlConfig); err != nil {
		log.Fatalln("mysql config failed: ", err)
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		mysqlConfig.User,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Dbname,
	)

	loggers := NewLogger(xlogger, logger.Config{
		SlowThreshold:             time.Second,
		Colorful:                  mysqlConfig.Colorful,
		IgnoreRecordNotFoundError: true,
		LogLevel:                  logger.LogLevel(mysqlConfig.LogLevel),
	})

	conn, err := gorm.Open(mysql.New(mysql.Config{DriverName: "mysql", DSN: dsn}), &gorm.Config{
		Logger: loggers,
	})
	if err != nil {
		log.Fatalln("gorm open err:", err)
	}

	db, err := conn.DB()
	if err != nil {
		log.Fatalln("mysql connect failed: ", err)
	}

	db.SetMaxIdleConns(mysqlConfig.MaxIdleConn)
	db.SetMaxOpenConns(mysqlConfig.MaxOpenConn)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		log.Fatalln("mysql ping failed: ", err)
	}
	log.Println("mysql connect successfully")

	return conn
}
