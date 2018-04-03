package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/astaxie/beego/config"
	"github.com/jinzhu/gorm"
	"os"
)

type Mconfig struct {
	NoticeLog bool
	Debuglog  bool
	Errlog    bool
	DBlog     bool

	NoticeLogFile string
	DebugLogFile  string
	ErrLogFile    string
	DBLogFile     string

	dBTestHost     string
	dBTestName     string
	dBTestUsername string
	dBTestPassword string

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

	//LoginTestServerDomain         string
	//LoginProcServerDomain         string

	ApiToken string

	// 默认Response消息
	ConfigMyResponse map[int]string

	// 是否在Docker容器中
	INDOCKER string
}

var (
	MyConfig     Mconfig
	dbOrmDefault *gorm.DB
)

const (
	//公共响应码
	RESP_OK        = 10000
	RESP_ERR       = 10001
	RESP_PARAM_ERR = 10002
	RESP_TOKEN_ERR = 10003
	RESP_NO_ACCESS = 10004
)

//AccessToken 相关
const (
	ROLE_STUDENT = 1
	ROLE_TEACHER = 2

	PLATFORM_ANDROID = 1
	PLATFORM_WEB     = 2
	PLATFORM_WEBCHAT = 3
)

const (
	//提取一些ID出来方便操作，如果修改数据库，记得修改这里。（一般不会动数据库ID吧）
	GRADE_ID_ZK = 12 //中考的年级ID
	GRADE_ID_GS = 13 //高三的年级ID

	COURSE_ID_WZ  = 98 //高考文综课程ID
	COURSE_ID_LZ  = 99 //高考理综课程ID
	COURSE_ID_GSX = 43 //高中数学的课程ID
	COURSE_ID_WX  = 14 //高考文数的课程ID
	COURSE_ID_LX  = 13 //高考理数的课程ID

	COURSE_ID_XXYW = 30 //小学英语课程ID
	COURSE_ID_XXSS = 31 //小学数学课程ID
	COURSE_ID_XXYY = 32 //小学英语课程ID
)

//收集问题接口的 类型集合
const (
	Word_Incorrect              = 1
	Answer_Incorrect            = 2
	Solution_Incorrect          = 3
	Question_Overflow_Incorrect = 4
	Other_Incorrect             = 5
)

//资源
type ResourceInfo struct {
	F_resource_id    int64       `json:"F_resource_id"`
	F_title          string      `json:"F_title"`
	F_data           interface{} `json:"F_data"`
	F_type_detail    int         `json:"F_type_detail"`
	F_courseware_uri string      `json:"F_courseware_uri"`
}

type MResp struct {
	F_responseNo  int    `required:"true" description:"响应码"`
	F_responseMsg string `description:"响应码描述"`
}

func init() {
	DREAMENV := "DEV"
	appConf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		return
	}

	MyConfig = Mconfig{}
	if appConf != nil {
		MyConfig.INDOCKER = os.Getenv("INDOCKER")

		MyConfig.NoticeLog, _ = appConf.Bool(DREAMENV + "::noticeLog")
		MyConfig.Debuglog, _ = appConf.Bool(DREAMENV + "::debuglog")
		MyConfig.Errlog, _ = appConf.Bool(DREAMENV + "::errlog")
		MyConfig.DBlog, _ = appConf.Bool(DREAMENV + "::dblog")
		MyConfig.NoticeLogFile = appConf.String(DREAMENV + "::noticeLogFile")
		MyConfig.DebugLogFile = appConf.String(DREAMENV + "::debugLogFile")
		MyConfig.ErrLogFile = appConf.String(DREAMENV + "::errLogFile")
		MyConfig.DBLogFile = appConf.String(DREAMENV + "::dbLogFile")

		MyConfig.dBTestHost = appConf.String(DREAMENV + "::dBTestHost")
		MyConfig.dBTestName = appConf.String(DREAMENV + "::dBTestName")
		MyConfig.dBTestUsername = appConf.String(DREAMENV + "::dBTestUsername")
		MyConfig.dBTestPassword = appConf.String(DREAMENV + "::dBTestPassword")

		MyConfig.dBProcHost = appConf.String(DREAMENV + "::dBProcHost")
		MyConfig.dBProcName = appConf.String(DREAMENV + "::dBProcName")
		MyConfig.dBProcUsername = appConf.String(DREAMENV + "::dBProcUsername")
		MyConfig.dBProcPassword = appConf.String(DREAMENV + "::dBProcPassword")

		MyConfig.dBMaxIdle, _ = appConf.Int(DREAMENV + "::dBMaxIdle")
		MyConfig.dBMaxConn, _ = appConf.Int(DREAMENV + "::dBMaxConn")

		MyConfig.SnowTestFlakDomain = appConf.String(DREAMENV + "::snowTestFlakDomain")
		MyConfig.SnowTestFlakAuthUser = appConf.String(DREAMENV + "::snowTestFlakAuthUser")
		MyConfig.SnowTestFlakAuthUserSecurity = appConf.String(DREAMENV + "::snowTestFlakAuthUserSecurity")

		MyConfig.SnowProcFlakDomain = appConf.String(DREAMENV + "::snowProcFlakDomain")
		MyConfig.SnowProcFlakAuthUser = appConf.String(DREAMENV + "::snowProcFlakAuthUser")
		MyConfig.SnowProcFlakAuthUserSecurity = appConf.String(DREAMENV + "::snowProcFlakAuthUserSecurity")

		//MyConfig.LoginTestServerDomain = appConf.String(DREAMENV + "::loginTestServerDomain")
		//MyConfig.LoginProcServerDomain = appConf.String(DREAMENV + "::loginProcServerDomain")
		MyConfig.ApiToken = appConf.String(DREAMENV + "::apiToken")
	}
	getResponseConfig()
}

//获取config
func getResponseConfig() {
	MyConfig.ConfigMyResponse = make(map[int]string)
	MyConfig.ConfigMyResponse[RESP_OK] = "成功"
	MyConfig.ConfigMyResponse[RESP_ERR] = "失败,未知错误"
	MyConfig.ConfigMyResponse[RESP_PARAM_ERR] = "参数错误"
	MyConfig.ConfigMyResponse[RESP_TOKEN_ERR] = "token错误"
	MyConfig.ConfigMyResponse[RESP_NO_ACCESS] = "没有访问权限"
}

//获取对应的db对象
func GetDb() *gorm.DB {
	return dbOrmDefault
}

func isIndocker() bool {
	return len(MyConfig.INDOCKER) > 0
}
