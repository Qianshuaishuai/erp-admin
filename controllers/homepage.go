package controllers

import (
	"elite-admin/models"
)

type HomePageController struct {
	BaseController
}

func (self *HomePageController) List() {
	self.Data["pageTitle"] = "首页标签列表"
	self.Data["ApiCss"] = true
	self.display()
}

func (self *HomePageController) Table() {
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
		models.DeleteHomeShow(id)
	}

	result, count := models.GetHomePageTagListSimple(q, limit, page)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.HomeShow.ID
		row["name"] = v.HomeShow.Name
		row["icon"] = v.HomeShow.Icon
		row["url"] = v.HomeShow.URL
		if v.HomeShow.Type == 2 {
			row["type"] = "外链标签"
		} else {
			row["type"] = "内容标签"
			row["url"] = "无外链"
		}
		if v.HomeShow.Type == 1 {
			var content string
			for c := range v.HomeContents {
				if v.HomeContents[c].Type == 1 {
					content += "图片：" + v.HomeContents[c].Content + "<br>"
				} else if v.HomeContents[c].Type == 2 {
					content += "段落：" + v.HomeContents[c].Content + "<br>"
				}
			}
			row["content"] = content
		} else if v.HomeShow.Type == 2 {
			row["content"] = "无正文"
		}

		list[k] = row
	}
	self.ajaxList("", 0, count, list)
}
