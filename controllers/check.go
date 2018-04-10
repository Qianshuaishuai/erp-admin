package controllers

import (
	"github.com/astaxie/beego"
	"dreamEbagPaperAdmin/models"
	"encoding/json"
	"strings"
	"dreamEbagPaperAdmin/helper"
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

	CHECK_STATUS_FLAG = [3]string{
		"<span class='layui-badge layui-bg-orange'>待提交</span>",
		"<span class='layui-badge layui-bg-green'>已提交</span>",
		"<span class='layui-badge layui-bg-gray'>已撤销</span>",
	}
)

func (self *CheckController) List() {
	self.Data["pageTitle"] = "审核列表"
	self.Data["ApiCss"] = true

	self.Data["IsChecker"] = self.isChecker()

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
		row["status_flag_text"] = CHECK_STATUS_FLAG[v.StatusFlag]
		row["status_flag"] = v.StatusFlag
		list[k] = row
	}
	self.ajaxList("", 0, count, list)
}

func (self *CheckController) Detail() {
	self.Data["pageTitle"] = "审核列表"
	self.Data["ApiCss"] = true

	modifyId, _ := self.GetInt64("modify_id")

	detailString := models.FindDetailById(modifyId)
	var Data []models.HistoryDetail
	json.Unmarshal([]byte(detailString), &Data)

	self.Data["Detail"] = Data
	self.display("check/detail")
}

func (self *CheckController) Delete() {
	if self.isChecker() {
		handler(models.DeleteCheckDataIds, self)
	} else {
		self.ajaxMsg("你没有对应权限", -1)
	}
}

func (self *CheckController) Revert() {
	handler(models.RevertCheckDataIds, self)
}

func (self *CheckController) Commit() {
	if self.isChecker() {
		handler(models.CommitCheckDataIds, self)
	} else {
		self.ajaxMsg("你没有对应权限", -1)
	}
}

func handler(dataHandler func([]int64) error, self *CheckController) {
	modifyIdsStr := strings.TrimSpace(self.GetString("ids"))

	if len(modifyIdsStr) != 0 {
		modifyIdsStr := strings.TrimRight(modifyIdsStr, ",")

		modifyIds, err := helper.TransformStringToInt64Arr("[" + modifyIdsStr + "]")
		if err != nil {
			self.ajaxMsg(err.Error(), -1)
		}

		err = dataHandler(modifyIds)

		if err != nil {
			self.ajaxMsg(err.Error(), -1)
		}
		self.ajaxMsg("", 0)
	}
}
