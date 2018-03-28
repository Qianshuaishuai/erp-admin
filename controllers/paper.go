package controllers

import (
	"github.com/astaxie/beego"
	"dreamEbagPaperAdmin/models"
	"strings"
	"dreamEbagPaperAdmin/helper"
)

type PaperController struct {
	BaseController
}

func (self *PaperController) List() {
	self.Data["pageTitle"] = "试卷列表"
	self.Data["ApiCss"] = true
	self.display()
}

func (self *PaperController) Table() {
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

	sort, err := self.GetInt("sort")
	if err != nil {
		sort = 0
	}

	result, count := models.GetPaperListSimple(q, limit, page, sort)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["paper_id"] = v.PaperId
		row["paper_name"] = v.Name
		row["paper_type_name"] = v.PaperTypeName
		row["paper_update_time"] = beego.Date(v.Date, "Y-m-d H:i:s")
		list[k] = row
	}
	self.ajaxList("", 0, count, list)
}

func (self *PaperController) Detail() {
	self.Data["ApiCss"] = true
	self.Data["pageTitle"] = "试卷详情"
	id, _ := self.GetInt64("paper_id", 0)
	provinces := models.GetProvinces()

	paperTypes := models.GetAllPaperType()
	papers := models.GetPaper(id)

	//PaperType
	var difficulty = 0

	if papers.Difficulty < 2.0 {
		difficulty = 2
	} else if papers.Difficulty >= 2.0 && papers.Difficulty < 3.9 {
		difficulty = 3
	} else if papers.Difficulty >= 3.9 && papers.Difficulty < 5.0 {
		difficulty = 4
	} else if papers.Difficulty >= 5.0 && papers.Difficulty < 6.1 {
		difficulty = 5
	} else if papers.Difficulty >= 6.1 {
		difficulty = 6
	}

	self.Data["Difficulty"] = difficulty
	self.Data["typeList"] = paperTypes
	self.Data["PaperType"] = int(papers.PaperType)
	self.Data["Detail"] = papers
	self.Data["UpdateTime"] = beego.Date(papers.Date, "Y-m-d H:i:s")
	self.Data["ProvinceList"] = provinces

	provinceMap := make(map[uint]string)
	for i := range provinces {
		for j := range papers.Provinces {
			if provinces[i].ProvinceId == papers.Provinces[j].ProvinceId {
				provinceMap[provinces[i].ProvinceId] = "checked"
			}
		}
	}
	self.Data["ProvinceMap"] = provinceMap
	resIds, _ := helper.TransformStringToInt64Arr(papers.QuestionSet.QuestionIds)
	self.Data["QuestionLens"] = len(resIds)

	chapters := papers.QuestionSet.PaperQuestionSetChapters
	var q = 0
	ChapterResult := make(map[int][]int64)
	for i := range chapters {
		chapterQuestionCount := chapters[i].QuestionCount
		a, b := models.GetTheQuestionByQ(resIds, q, int(chapterQuestionCount))
		ChapterResult[i] = resIds[a:b]
		q = b
	}

	self.Data["ChapterResult"] = ChapterResult

	self.display()
}

func (self *PaperController) Edit() {
	Paper_id, _ := self.GetInt64("paper_id")
	if Paper_id != 0 {
		paper_name := strings.TrimSpace(self.GetString("paper_name"))
		paper_full_score, _ := self.GetInt("full_score", -100)
		paper_type, _ := self.GetInt("paper_type", -100)
		difficulty, _ := self.GetFloat("difficulty", -100)
		provinces := strings.TrimSpace(self.GetString("province"))

		if err := models.UpdatePaper(Paper_id, paper_name, paper_full_score, paper_type, difficulty, provinces); err != nil {
			self.ajaxMsg(err.Error(), -1)
		}
		self.ajaxMsg("", 0)
	}
}
