package controllers

import "elite-admin/models"

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

	result, count := models.GetIndustryTagListSimple(q, limit, page)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name

		list[k] = row
	}
	self.ajaxList("", 0, count, list)
}

func (self *IndustryController) Add() {
	self.Data["pageTitle"] = "添加标签"
	self.Data["ApiCss"] = true
	self.display()
}

func (self *IndustryController) AddTag() {
	name := self.GetString("name")

	err := models.AddIndustryTag(name)

	if err != nil {
		self.ajaxMsg("添加失败 :"+err.Error(), -1)
	}

	self.ajaxMsg("success", 0)
}
