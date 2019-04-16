package controllers

import (
	"elite-admin/models"
	"strconv"

	"github.com/astaxie/beego"
)

type ConnectionController struct {
	BaseController
}

func (self *ConnectionController) List() {
	self.Data["pageTitle"] = "专家列表"
	self.Data["ApiCss"] = true
	self.display()
}

func (self *ConnectionController) Table() {
	//列表
	page, err := self.GetInt("page")
	if err != nil {
		page = 1
	}
	limit, err := self.GetInt("limit")
	if err != nil {
		limit = 30
	}

	// q 查询条件
	q := self.GetString("q")

	phone, _ := self.GetInt("id")
	status, _ := self.GetInt("status")

	if phone > 0 && status > 0 {
		models.ChangeConnectionStatus(phone, status)
	}

	//看看q能不能变成ID
	paperId, err := strconv.ParseInt(q, 10, 64)

	if err == nil {
		q = ""
	}

	sort, err := self.GetInt("sort")
	if err != nil {
		sort = 0
	}

	result, count := models.GetConnectionListSimple(q, limit, page, sort, paperId)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["phone"] = v.Phone
		row["good"] = v.Good
		row["look"] = v.Look
		row["card"] = v.Card
		row["time"] = beego.Date(v.Time, "Y-m-d H:i:s")
		if v.Status == 0 {
			row["status"] = "未审核"
		} else if v.Status == 1 {
			row["status"] = "审核通过"
		} else if v.Status == 2 {
			row["status"] = "审核不通过"
		}
		list[k] = row
	}
	self.ajaxList("", 0, count, list)
}

func (self *ConnectionController) Add() {
	personTags, _ := models.GetPersonTagListSimple("", 1000, 1)
	self.Data["pageTitle"] = "添加新专家"
	self.Data["ApiCss"] = true
	self.Data["PersonTagList"] = personTags
	self.display()
}

func (self *ConnectionController) AddConnection() {
	username := self.GetString("username")
	job := self.GetString("job")
	position := self.GetString("position")
	profess := self.GetString("profess")
	agency := self.GetString("agency")
	address := self.GetString("address")
	introduce := self.GetString("introduce")
	achieve := self.GetString("achieve")
	school := self.GetString("school")
	phone, _ := self.GetInt("phone")
	look, _ := self.GetInt("look")
	good, _ := self.GetInt("good")
	tags := self.GetString("tags")

	iFile, iHandler, _ := self.GetFile("iconFile")
	cFile, cHandler, _ := self.GetFile("cardFile")

	iImageURL, _ := models.UploadFile(models.TYPE_CONNECTION_ICON_ID, iHandler.Filename, iFile)
	cImageURL, _ := models.UploadFile(models.TYPE_CONNECTION_CARD_ID, cHandler.Filename, cFile)

	err := models.AddConnection(phone, look, good, username, job, position, profess, agency, address, introduce, achieve, school, iImageURL, cImageURL, tags)

	if err != nil {
		self.ajaxMsg("添加失败 :"+err.Error(), -1)
	}

	self.ajaxMsg("success", 0)
}
