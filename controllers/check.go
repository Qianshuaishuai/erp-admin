package controllers

import (
	"github.com/astaxie/beego"
	"dreamEbagPaperAdmin/models"
)

type CheckController struct {
	BaseController
}

var (
	DATA_TYPA_MAP = map[int]string{
		1: "试卷",
		2: "章节",
		3: "小题",
		4: "大题",
	}

	DATA_OPERATE_MAP = map[int]string{
		101: "编辑",
		201: "增加",
		301: "删除",
	}

	COMMIT_FLAG = [2]string{
		"<span class='layui-badge layui-bg-yellow'>待提交</span>",
		"<span class='layui-badge layui-bg-green'>已提交</span>",
	}
)

func (self *CheckController) List() {
	self.Data["pageTitle"] = "审核列表"
	self.Data["ApiCss"] = true
	self.display()
}

func (self *CheckController) Table() {
	//列表
	page, err := self.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := self.GetInt("limit")
	if err != nil {
		limit = 10
	}

	result, count := models.GetCheckList(limit, page)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["modify_id"] = v.ModifyId
		row["modify_date"] = beego.Date(v.ModifyDate, "Y-m-d H:i:s")
		row["modify_admin"] = v.ModifyAdmin
		row["data_id"] = v.DataId
		row["data_type"] = DATA_TYPA_MAP[v.DataType]
		row["data_operate"] = DATA_OPERATE_MAP[v.DataOperate]
		row["commit_flag"] = COMMIT_FLAG[v.CommitFlag]
		list[k] = row
	}
	self.ajaxList("", 0, count, list)
}
