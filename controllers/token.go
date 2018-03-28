package controllers

import (
	"dreamEbagPapers/models"
	"sync"
)

type TokenController struct {
	BaseController
}

// 全局声明此Map，注意，其它地方只能读取值，不要写入！
var TeacherTokenMap sync.Map

//@Title 获取保存凭证
//@Description 根据老师ID获取保存资源需要提供的F_token，过期时间：1小时。
// @Param	F_teacher_id    	form	string		true	用户ID
// @Param	F_accesstoken		form	string		true	ACCESSTOKEN
//@OutParam F_responseNo int true 业务错误码（可能值：10000,10001,10002）
//@OutParam F_responseMsg string false 业务错误描述
//@OutParam F_token string true 保存凭证
//@Router v1/token [get]
func (c *TokenController) GetSaveToken() {
	datas := map[string]interface{}{"F_responseNo": models.RESP_OK}
	F_teacher_id := c.Ctx.Input.Query("F_teacher_id")

	token := models.GetSaveToken()
	datas["F_token"] = token
	TeacherTokenMap.Store(F_teacher_id, token)

	c.jsonEcho(datas)
}
