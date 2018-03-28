package controllers

import (
	"dreamEbagPaperAdmin/models"
	"fmt"
)

type EnvController struct {
	BaseController
}

func (self *EnvController) ChangeToTest() {
	datas := map[string]interface{}{"F_responseNo": models.RESP_OK}
	fmt.Println("Change To Test")
	models.ChangeDBToTest()
	self.jsonEcho(datas)
}

func (self *EnvController) ChangeToProc() {
	datas := map[string]interface{}{"F_responseNo": models.RESP_OK}
	fmt.Println("Change To Proc")
	models.ChangeDBToProc()
	self.jsonEcho(datas)
}
