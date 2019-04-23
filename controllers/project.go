package controllers

import (
	"elite-admin/models"
	"strconv"

	"github.com/astaxie/beego"
)

type ProjectController struct {
	BaseController
}

func (self *ProjectController) List() {
	self.Data["pageTitle"] = "项目列表"
	self.Data["ApiCss"] = true
	self.display()
}

func (self *ProjectController) Table() {
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

	id, _ := self.GetInt("id")
	status, _ := self.GetInt("status")

	if id > 0 && status > 0 {
		models.ChangeProjectStatus(id, status)
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

	result, count := models.GetProjectListSimple(q, limit, page, sort, paperId)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name
		row["phone"] = v.Phone
		row["icon"] = v.Icon
		row["background"] = v.Background
		row["type"] = v.Type
		row["address"] = v.Address
		row["addtip"] = v.AddTip
		row["money"] = v.Money
		row["agency"] = v.Agency
		row["introduce"] = v.Introduce
		row["idcard"] = v.IDCard
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

func (self *ProjectController) Add() {
	industrys, _ := models.GetIndustryTagListSimple("", 1000, 1)
	self.Data["pageTitle"] = "添加新项目"
	self.Data["ApiCss"] = true
	self.Data["IndustryList"] = industrys
	self.display()
}

func (self *ProjectController) Detail() {
	id, _ := self.GetInt64("id", 0)
	data, _ := models.GetProjectDetail(id)
	self.Data["pageTitle"] = "项目详情"
	self.Data["ApiCss"] = true
	self.Data["Project"] = data.Project
	self.Data["Industrys"] = data.Industrys
	self.Data["Time"] = beego.Date(data.Project.Time, "Y-m-d H:i:s")

	if data.Project.Status == 0 {
		self.Data["Status"] = "未审核"
	} else if data.Project.Status == 1 {
		self.Data["Status"] = "审核通过"
	} else {
		self.Data["Status"] = "审核不通过"
	}

	self.display()
}

func (self *ProjectController) AddProject() {
	name := self.GetString("name")
	typeName := self.GetString("type")
	address := self.GetString("address")
	money := self.GetString("money")
	agency := self.GetString("agency")
	introduce := self.GetString("introduce")
	addtip := self.GetString("addtip")
	idcard := self.GetString("idcard")
	phone, _ := self.GetInt("phone")
	industrys := self.GetString("industrys")

	cFile, cHandler, _ := self.GetFile("cardFile")
	bFile, bHandler, _ := self.GetFile("backgroundFile")

	var cImageURL, bImageURL string

	if cFile == nil {
		cImageURL = ""
	} else {
		cImageURL, _ = models.UploadFile(models.TYPE_PROJECT_CARD_ID, cHandler.Filename, cFile)
	}

	if bFile == nil {
		bImageURL = ""
	} else {
		bImageURL, _ = models.UploadFile(models.TYPE_PROJECT_BACKGROUND_ID, bHandler.Filename, bFile)
	}

	err := models.AddProject(name, typeName, address, money, agency, introduce, addtip, idcard, phone, cImageURL, bImageURL, industrys)

	if err != nil {
		self.ajaxMsg("添加失败 :"+err.Error(), -1)
	}

	self.ajaxMsg("success", 0)
}
