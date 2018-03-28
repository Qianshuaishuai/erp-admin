package controllers

import (
	"github.com/astaxie/beego"
	"dreamEbagPaperAdmin/models"
)

type ErrorController struct {
	beego.Controller
}

func init() {
}

//json echo
func (u0 *ErrorController) jsonEcho(datas map[string]interface{}, u *ErrorController) {
	u.Ctx.Output.ContentType("application/json; charset=utf-8")
	u.Data["json"] = datas
	u.ServeJSON()
}

func (u *ErrorController) Error401() {
	//ini return
	datas := map[string]interface{}{"F_responseNo": 401}
	//log
	u.errerLog("401")
	//return
	u.jsonEcho(datas, u)
}

func (u *ErrorController) Error403() {
	//ini return
	datas := map[string]interface{}{"F_responseNo": 403}
	//log
	u.errerLog("403")
	//return
	u.jsonEcho(datas, u)
}

func (u *ErrorController) Error404() {
	//ini return
	datas := map[string]interface{}{"F_responseNo": 404}
	//log
	u.errerLog("404")
	//return
	u.jsonEcho(datas, u)
}

func (u *ErrorController) Error500() {
	//ini return
	datas := map[string]interface{}{"F_responseNo": 500}
	//log
	u.errerLog("500")
	//return
	u.jsonEcho(datas, u)
}

func (u *ErrorController) Error503() {
	//ini return
	datas := map[string]interface{}{"F_responseNo": 503}
	//log
	u.errerLog("503")
	//return
	u.jsonEcho(datas, u)
}

func (u *ErrorController) errerLog(code string) {
	var logObj models.MLog
	logObj.LogRequestErrCode(u.Ctx, code)
}
