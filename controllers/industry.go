package controllers

import "erp-admin/models"

type IndustryController struct {
	BaseController
}

func (self *IndustryController) List() {
	self.Data["pageTitle"] = "个人标签列表"
	self.Data["ApiCss"] = true
	self.display()
}

func (self *IndustryController) Table() {
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
	if id > 0 {
		models.DeleteIndustryTag(id)
	}

	sort, err := self.GetInt("sort")
	if err != nil {
		sort = 0
	}

	result, count := models.GetIndustryTagListSimple(q, limit, page, sort)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name
		row["plain"] = v.Plain

		list[k] = row
	}
	self.ajaxList("", 0, count, list)
}

func (self *IndustryController) Add() {
	self.Data["pageTitle"] = "添加标签"
	self.Data["ApiCss"] = true
	self.display()
}

func (self *IndustryController) Edit() {
	id, _ := self.GetInt64("id", 0)
	data, _ := models.GetIndustryDetail(id)
	self.Data["pageTitle"] = "编辑行业标签"
	self.Data["ApiCss"] = true
	self.Data["IndustryTag"] = data

	self.display()
}

func (self *IndustryController) EditTag() {
	id, _ := self.GetInt("id")
	index, _ := self.GetInt("index")
	name := self.GetString("name")

	err := models.EditIndustryDetail(id, index, name)

	if err != nil {
		self.ajaxMsg("编辑失败 :"+err.Error(), -1)
	}

	self.ajaxMsg("success", 0)
}

func (self *IndustryController) AddTag() {
	name := self.GetString("name")

	err := models.AddIndustryTag(name)

	if err != nil {
		self.ajaxMsg("添加失败 :"+err.Error(), -1)
	}

	self.ajaxMsg("success", 0)
}
