package controllers

import (
	"dreamEbagPaperAdmin/models"
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

	if err == nil {
		q = ""
	}

	result, count := models.GetPersonTagListSimple(q, limit, page)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.ID
		row["name"] = v.Name
		row["icon"] = v.Icon

		list[k] = row
	}
	self.ajaxList("", 0, count, list)
}

func (self *PersonController) Add() {
	self.Data["pageTitle"] = "添加标签"
	self.Data["ApiCss"] = true
	self.display()
}

func (self *PersonController) AddTag() {
	// name := self.GetString("name")
	// file, fileHandler, _ := self.GetFile("file")

	// imageUrl, _ := models.UploadFile(models.TYPE_USER_HEAD_ID, fileHandler.Filename, file)

	// err := models.AddPersonTag(name, imageUrl)

	// if err != nil {
	// 	self.ajaxMsg("添加标签失败 :"+err.Error(), -1)
	// }
	self.ajaxMsg("success", 0)
}
