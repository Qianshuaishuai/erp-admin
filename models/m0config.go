package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"os"

	"github.com/HYY-yu/LogLib"
	"github.com/astaxie/beego/config"
	"github.com/jinzhu/gorm"
)

type Mconfig struct {
	LogLevel int

	LogFile string

	dBTestHost     string
	dBTestName     string
	dBTestUsername string
	dBTestPassword string

	dBEliteHost     string
	dBEliteName     string
	dBEliteUsername string
	dBElitePassword string

	dBProcHost     string
	dBProcName     string
	dBProcUsername string
	dBProcPassword string

	dBMaxIdle int
	dBMaxConn int

	SnowTestFlakDomain           string
	SnowTestFlakAuthUser         string
	SnowTestFlakAuthUserSecurity string

	SnowProcFlakDomain           string
	SnowProcFlakAuthUser         string
	SnowProcFlakAuthUserSecurity string

	ApiToken string

	// 是否在Docker容器中
	INDOCKER string

	//错误处理
	RecoverPanic bool

	QiniuAccessKey string
	QiniuSecretKey string
	QiniuBucket    string
	QiniuBaseURL   string
}

var (
	MyConfig     Mconfig
	dbOrmDefault *gorm.DB
	dbOrmElite   *gorm.DB
)

const (
	//公共响应码
	RESP_OK        = 10000
	RESP_ERR       = 10001
	RESP_PARAM_ERR = 10002
	RESP_TOKEN_ERR = 10003
	RESP_NO_ACCESS = 10004
)

const (
	ADMIN_SUPER   = -1
	ADMIN_DATAER  = 1
	ADMIN_CHECKER = 2
)

func init() {
	DREAMENV := "DEV"
	appConf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		return
	}

	MyConfig = Mconfig{}
	if appConf != nil {
		MyConfig.INDOCKER = os.Getenv("INDOCKER")

		levelStr := appConf.String(DREAMENV + "::LogLevel")
		switch levelStr {
		case "DEBUG":
			MyConfig.LogLevel = loglib.LevelDebug
		case "INFO":
			MyConfig.LogLevel = loglib.LevelInfo
		case "ERROR":
			MyConfig.LogLevel = loglib.LevelError
		}

		MyConfig.LogFile = appConf.String(DREAMENV + "::LogFile")

		MyConfig.dBTestHost = appConf.String(DREAMENV + "::dBTestHost")
		MyConfig.dBTestName = appConf.String(DREAMENV + "::dBTestName")
		MyConfig.dBTestUsername = appConf.String(DREAMENV + "::dBTestUsername")
		MyConfig.dBTestPassword = appConf.String(DREAMENV + "::dBTestPassword")

		MyConfig.dBProcHost = appConf.String(DREAMENV + "::dBProcHost")
		MyConfig.dBProcName = appConf.String(DREAMENV + "::dBProcName")
		MyConfig.dBProcUsername = appConf.String(DREAMENV + "::dBProcUsername")
		MyConfig.dBProcPassword = appConf.String(DREAMENV + "::dBProcPassword")

		MyConfig.dBEliteHost = appConf.String(DREAMENV + "::dBEliteHost")
		MyConfig.dBEliteName = appConf.String(DREAMENV + "::dBEliteName")
		MyConfig.dBEliteUsername = appConf.String(DREAMENV + "::dBEliteUsername")
		MyConfig.dBElitePassword = appConf.String(DREAMENV + "::dBElitePassword")

		MyConfig.dBMaxIdle, _ = appConf.Int(DREAMENV + "::dBMaxIdle")
		MyConfig.dBMaxConn, _ = appConf.Int(DREAMENV + "::dBMaxConn")

		MyConfig.SnowTestFlakDomain = appConf.String(DREAMENV + "::snowTestFlakDomain")
		MyConfig.SnowTestFlakAuthUser = appConf.String(DREAMENV + "::snowTestFlakAuthUser")
		MyConfig.SnowTestFlakAuthUserSecurity = appConf.String(DREAMENV + "::snowTestFlakAuthUserSecurity")

		MyConfig.SnowProcFlakDomain = appConf.String(DREAMENV + "::snowProcFlakDomain")
		MyConfig.SnowProcFlakAuthUser = appConf.String(DREAMENV + "::snowProcFlakAuthUser")
		MyConfig.SnowProcFlakAuthUserSecurity = appConf.String(DREAMENV + "::snowProcFlakAuthUserSecurity")

		MyConfig.ApiToken = appConf.String(DREAMENV + "::apiToken")

		MyConfig.QiniuAccessKey = appConf.String(DREAMENV + "::accessKey")
		MyConfig.QiniuSecretKey = appConf.String(DREAMENV + "::secretKey")
		MyConfig.QiniuBucket = appConf.String(DREAMENV + "::bucket")
		MyConfig.QiniuBaseURL = appConf.String(DREAMENV + "::qiniuBaseUrl")
	}
	initLog()
}

func initLog() {
	//初始化日志模块
	if Indocker() {
		loglib.InitLogger(loglib.LogConfig{LogTo: loglib.ConsoleLogs, LogLevel: MyConfig.LogLevel, LogPretty: false})
	} else {
		loglib.InitLogger(loglib.LogConfig{LogTo: loglib.FileLogs, LogPath: MyConfig.LogFile, LogLevel: MyConfig.LogLevel, LogPretty: true})
	}
}

//获取对应的db对象
func GetDb() *gorm.DB {
	return dbOrmDefault
}

//获取新的db对象
func GetEliteDb() *gorm.DB {
	return dbOrmElite
}

func Indocker() bool {
	return len(MyConfig.INDOCKER) > 0
}
