package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"strings"
	"dreamEbagPaperAdmin/helper"
)

var (
	noticeLog *logs.BeeLogger
	debugLog  *logs.BeeLogger
	errLog    *logs.BeeLogger
)

var (
	defaultLogger MLog
	consoleLog    *logs.BeeLogger
)

type MLog struct {
}

//打Logger请使用GetLogger()方法。
func GetLogger() *MLog {
	return &defaultLogger
}

func logPrintDebug(format string, v ...interface{}) {
	if MyConfig.Debuglog {
		if isIndocker() {
			consoleLog.Debug(format, v)
		} else {
			debugLog.Info(format, v)
		}
	}
}

func logPrintInfo(format string, v ...interface{}) {
	if MyConfig.NoticeLog {
		if isIndocker() {
			consoleLog.Info(format, v)
		} else {
			noticeLog.Info(format, v)
		}
	}
}

func logPrintWarn(format string, v ...interface{}) {
	if MyConfig.Errlog {
		if isIndocker() {
			consoleLog.Warn(format, v)
		} else {
			errLog.Info(format, v)
		}
	}
}

func init() {
	logs.EnableFuncCallDepth(false)

	if isIndocker() {
		//使用控制台日志
		consoleLog = logs.NewLogger()
		consoleLog.SetLogger(logs.AdapterConsole)
		consoleLog.Async()

		//警告日志
		if MyConfig.Errlog {
			consoleLog.SetLevel(logs.LevelWarn)
		}

		//消息日志
		if MyConfig.NoticeLog {
			consoleLog.SetLevel(logs.LevelInfo)
		}

		//所有日志 (调试日志以上)
		if MyConfig.Debuglog {
			consoleLog.SetLevel(logs.LevelTrace)
		}
	} else {
		//使用文件日志
		noticeLog = logs.NewLogger(10000)
		str, _ := json.Marshal(map[string]interface{}{"filename": MyConfig.NoticeLogFile, "perm": "0644"})
		noticeLog.SetLogger("file", string(str))
		noticeLog.Async()

		debugLog = logs.NewLogger(10000)
		str, _ = json.Marshal(map[string]interface{}{"filename": MyConfig.DebugLogFile, "perm": "0644"})
		debugLog.SetLogger("file", string(str))
		debugLog.Async()

		errLog = logs.NewLogger(10000)
		str, _ = json.Marshal(map[string]interface{}{"filename": MyConfig.ErrLogFile, "perm": "0644"})
		errLog.SetLogger("file", string(str))
		errLog.Async()
	}
}

//继承自elastic.Logger 记录ES错误日志
func (u0 *MLog) Printf(format string, v ...interface{}) {
	uniqueLogFlag := helper.GetGuid()
	str := " ElasticSearch - uniqueLogFlag:" + uniqueLogFlag + " - " + fmt.Sprintf(format, v)
	logPrintWarn(str)
}

//记录请求(log等级)
func (u0 *MLog) LogRequest(u *context.Context, uniqueLogFlag string) {
	u.Request.ParseForm()
	str := " Request - uniqueLogFlag:" + uniqueLogFlag + " \n " + GetDisplayString("IP", u.Input.IP(), "Scheme", u.Request.Proto, "Uri", u.Request.RequestURI, "Method", u.Request.Method, "Post", u.Request.Form, "Header", u.Request.Header)
	logPrintDebug(str)
}

//记录返回(log等级)
func (u0 *MLog) LogEcho(datas interface{}, uniqueLogFlag string) {
	b, err := json.Marshal(datas)
	if err == nil {
		str := " Echo - uniqueLogFlag:" + uniqueLogFlag + " \n " + GetDisplayString("datas", string(b))
		logPrintDebug(str)
	}
}

func (u0 *MLog) LogRequestErrCode(u *context.Context, code string) {
	u.Request.ParseForm()
	str := " Request error - code:" + code + " \n " + GetDisplayString("IP", u.Input.IP(), "Scheme", u.Request.Proto, "Uri", u.Request.RequestURI, "Method", u.Request.Method, "Post", u.Request.Form, "Header", u.Request.Header)
	logPrintWarn("%s", str)
}

//debug(debug等级)
func (u0 *MLog) LogDebug(data interface{}, uniqueLogFlag string) {
	str := " Debug - uniqueLogFlag:" + uniqueLogFlag + " \n " + GetDisplayString("debug:", data)
	logPrintDebug("%s", str)
}

//记录Info(log等级)
func (u0 *MLog) LogInfo(data interface{}, uniqueLogFlag string) {
	str := " Info - uniqueLogFlag:" + uniqueLogFlag + " \n " + GetDisplayString("info:", data)
	logPrintInfo("%s", str)
}

//记录err(warn等级)
func (u0 *MLog) LogErr(err interface{}, uniqueLogFlag string, descriptions ...string) {
	var description string
	if len(descriptions) > 0 {
		description = strings.Join(descriptions, "-")
	}
	str := " Err - uniqueLogFlag:" + uniqueLogFlag + " \n " + description + " : " + GetDisplayString("err", err)
	logPrintWarn("%s", str)
}

//snowflak curl request(log等级)
func (u0 *MLog) LogSnowflakCurlRequest(uri string, method string, data interface{}) {
	str := GetDisplayString("Snowflak curl uri:", uri, "Snowflak curl method:", method, "Snowflak curl data:", data)
	logPrintDebug("%s", str)
}

//snowflak curl response(log等级)
func (u0 *MLog) LogSnowflakCurlResponse(body interface{}, header interface{}, status interface{}) {
	str := GetDisplayString("Snowflak curl response body:", body, "Snowflak curl response header:", header, "Snowflak curl reponse status:", status)
	logPrintDebug("%s", str)
}

//other curl request(log等级)
func (u0 *MLog) LogOtherCurlRequest(uniqueLogFlag string, description string, uri string, method string, data interface{}) {
	str := "\nRquest,uniqueLogFlag:" + uniqueLogFlag + ", " + description + "\n" + GetDisplayString("curl uri:", uri, "curl method:", method, "curl request data:", data)
	logPrintDebug("%s", str)
}

//other curl response(log等级)
func (u0 *MLog) LogOtherCurlResponse(uniqueLogFlag string, description string, body interface{}, header interface{}, status interface{}) {
	str := "\nRespose,uniqueLogFlag:" + uniqueLogFlag + ", " + description + "\n" + GetDisplayString("curl response body:", body, "curl response header:", header, "curl reponse status:", status)
	logPrintDebug("%s", str)
}

//other curl err(log等级)
func (u0 *MLog) LogOtherCurlResponseErr(uniqueLogFlag string, description string, err interface{}) {
	str := "\nerr,uniqueLogFlag:" + uniqueLogFlag + "\n" + description + "   :    " + GetDisplayString("err", err)
	logPrintDebug("%s", str)
}
