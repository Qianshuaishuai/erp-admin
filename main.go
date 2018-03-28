package main

import (
	"github.com/astaxie/beego"
	"fmt"
	"strconv"
	"runtime"
	"dreamEbagPaperAdmin/controllers"
	"dreamEbagPaperAdmin/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/context"

	_ "dreamEbagPaperAdmin/routers"
)

func main() {
	//连接数据库
	models.InitGorm()
	db := models.GetDb()
	defer db.Close()

	//就问能不能ping通
	errPing := db.DB().Ping()
	if errPing != nil {
		fmt.Println("can't connect db", errPing.Error())
		return
	}

	//如果服务器Panic ，返回500错，而不是错误信息。
	//beego.BConfig.RecoverFunc = recoverFuncForServer
	//beego.BConfig.RecoverPanic = false
	//beego.ErrorController(&controllers.ErrorController{})
	beego.Run()
}

func recoverFuncForServer(ctx *context.Context) {
	if err := recover(); err != nil {
		if err == beego.ErrAbort {
			return
		}
		if !beego.BConfig.RecoverPanic {
			panic(err)
		}
		if beego.BConfig.EnableErrorsShow {
			if _, ok := beego.ErrorMaps[fmt.Sprint(err)]; ok {
				x, e := strconv.ParseUint(fmt.Sprint(err), 10, 64)
				if e == nil {
					beego.Exception(x, ctx)
					return
				}
			}
			//panic就抛一个500错
			beego.Exception(500, ctx)
		}

		var stack string
		models.GetLogger().LogInfo("the request url is "+ctx.Input.URL(), "")
		models.GetLogger().LogErr(err, "Handler crashed with error")
		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			logs.Critical(fmt.Sprintf("%s:%d", file, line))
			stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
		}
	}
}
