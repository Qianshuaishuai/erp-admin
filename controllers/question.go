package controllers

import (
	"dreamEbagPaperAdmin/models"
	"dreamEbagPaperAdmin/helper"
)

type QuestionController struct {
	BaseController
}

func (self *QuestionController) Search() {
	self.Data["pageTitle"] = "搜索试题"
	self.Data["ApiCss"] = true

	self.display()
}

func (self *QuestionController) Detail() {
	self.Data["pageTitle"] = "试题详情"
	self.Data["ApiCss"] = true

	questionId, _ := self.GetInt64("question_id", 0)
	isBig ,data :=models.GetQuestion(questionId)

	self.Data["IsBig"] = isBig
	self.Data["Data"] = data

	if isBig{
		bigQuestion := data.(models.BigQuestion)
		smallQuestionIds,_ := helper.TransformStringToInt64Arr(bigQuestion.BigQuestionIds)
		self.Data["Questions"] = smallQuestionIds
	}else{
		smallQuestion := data.(models.SmallQuestion)
		tempMap := make(map[int]string)
		tempMap[smallQuestion.RealType] = models.SmallQuestionType[smallQuestion.RealType]
		if smallQuestion.RealType == models.OBJECTIVELY_BLANK || smallQuestion.RealType == models.SUBJECTIVITY_BLANK {
			tempMap[models.OBJECTIVELY_BLANK] = "客观填空题"
			tempMap[models.SUBJECTIVITY_BLANK] = "主观填空题"
		}
		self.Data["QuestionTypeMap"] = tempMap
	}

	self.display()
}
