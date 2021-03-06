package controllers

import (
	"erp-admin/models"
)

type PersonController struct {
	BaseController
}

func (self *PersonController) List() {
	self.Data["pageTitle"] = "个人标签列表"
	self.Data["ApiCss"] = true
	self.display()
}

func (self *PersonController) Table() {
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
		models.DeletePersonTag(id)
	}

	sort, err := self.GetInt("sort")
	if err != nil {
		sort = 0
	}

	result, count := models.GetPersonTagListSimple(q, limit, page, sort)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name
		row["icon"] = v.Icon
		row["plain"] = v.Plain

		list[k] = row
	}
	self.ajaxList("", 0, count, list)
}

func (self *PersonController) Add() {
	self.Data["pageTitle"] = "添加标签"
	self.Data["ApiCss"] = true
	self.display()
}

func (self *PersonController) Edit() {
	id, _ := self.GetInt64("id", 0)
	data, _ := models.GetPersonDetail(id)
	self.Data["pageTitle"] = "编辑个人标签"
	self.Data["ApiCss"] = true
	self.Data["PersonTag"] = data

	self.display()
}

func (self *PersonController) EditTag() {
	id, _ := self.GetInt("id")
	index, _ := self.GetInt("index")
	name := self.GetString("name")
	imageURL := self.GetString("image")

	if imageURL == "" {
		file, fileHandler, _ := self.GetFile("file")
		imageURL, _ = models.UploadFile(models.TYPE_ADD_PERSON_TAG, fileHandler.Filename, file)
	}

	err := models.EditPersonDetail(id, index, name, imageURL)

	if err != nil {
		self.ajaxMsg("编辑失败 :"+err.Error(), -1)
	}

	self.ajaxMsg("success", 0)
}

func (self *PersonController) AddTag() {
	name := self.GetString("name")
	file, fileHandler, _ := self.GetFile("file")

	imageURL, _ := models.UploadFile(models.TYPE_ADD_PERSON_TAG, fileHandler.Filename, file)

	err := models.AddPersonTag(name, imageURL)

	if err != nil {
		self.ajaxMsg("添加失败 :"+err.Error(), -1)
	}

	self.ajaxMsg("success", 0)
}
