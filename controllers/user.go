package controllers

import (
	"dreamEbagPaperAdmin/helper"
	"dreamEbagPaperAdmin/models"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

type UserController struct {
	BaseController
}

func (self *UserController) List() {
	self.Data["pageTitle"] = "用户列表"
	self.Data["ApiCss"] = true
	self.display()
}

func (self *UserController) Table() {
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

	//看看q能不能变成ID
	paperId, err := strconv.ParseInt(q, 10, 64)

	if err == nil {
		q = ""
	}

	sort, err := self.GetInt("sort")
	if err != nil {
		sort = 0
	}

	result, count := models.GetUserListSimple(q, limit, page, sort, paperId)
	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["phone"] = v.Phone
		row["username"] = v.Username
		row["address"] = v.Address
		row["school"] = v.School
		row["register"] = beego.Date(v.Register, "Y-m-d H:i:s")
		list[k] = row
	}
	self.ajaxList("", 0, count, list)
}

func (self *UserController) Detail() {
	self.Data["ApiCss"] = true
	self.Data["pageTitle"] = "试卷详情"
	id, _ := self.GetInt64("paper_id", 0)
	provinces := models.GetProvinces()
	paperTypes := models.GetAllPaperType()
	paper := models.GetPaper(id)

	//PaperType
	var difficulty = 0

	if paper.Difficulty < 2.0 {
		difficulty = 2
	} else if paper.Difficulty >= 2.0 && paper.Difficulty < 3.9 {
		difficulty = 3
	} else if paper.Difficulty >= 3.9 && paper.Difficulty < 5.0 {
		difficulty = 4
	} else if paper.Difficulty >= 5.0 && paper.Difficulty < 6.1 {
		difficulty = 5
	} else if paper.Difficulty >= 6.1 {
		difficulty = 6
	}

	self.Data["Difficulty"] = difficulty
	self.Data["typeList"] = paperTypes
	self.Data["Detail"] = paper
	self.Data["UpdateTime"] = beego.Date(paper.Date, "Y-m-d H:i:s")
	self.Data["ProvinceList"] = provinces

	provinceMap := make(map[uint]string)
	for i := range provinces {
		for j := range paper.Provinces {
			if provinces[i].ProvinceId == paper.Provinces[j].ProvinceId {
				provinceMap[provinces[i].ProvinceId] = "checked"
			}
		}
	}
	self.Data["ProvinceMap"] = provinceMap
	resIds, _ := helper.TransformStringToInt64Arr(paper.QuestionSet.QuestionIds)
	self.Data["QuestionLens"] = len(resIds)

	chapters := paper.QuestionSet.PaperQuestionSetChapters
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

func (self *UserController) Edit() {
	paperId, _ := self.GetInt64("paper_id")
	if paperId != 0 {
		paperName := strings.TrimSpace(self.GetString("paper_name"))
		paperFullScore, _ := self.GetInt("full_score", -100)
		paperType, _ := self.GetInt("paper_type", -100)
		difficulty, _ := self.GetFloat("difficulty", -100)
		provinces := strings.TrimSpace(self.GetString("province"))

		if err := models.UpdatePaper(self.user, paperId, paperName, paperFullScore, paperType, difficulty, provinces); err != nil {
			self.ajaxMsg(err.Error(), -1)
		}
		self.ajaxMsg("", 0)
	}
}
