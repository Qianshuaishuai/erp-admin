package controllers

import (
	"erp-admin/helper"
	"erp-admin/models"
	"strconv"
	"strings"
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
	questionQ := self.GetString("question_q")
	isBig, data := models.GetQuestion(questionId, questionQ)

	if isBig {
		bigQuestion := data.(models.BigQuestion)
		smallQuestionIds, _ := helper.TransformStringToInt64Arr(bigQuestion.BigQuestionIds)
		self.Data["Questions"] = smallQuestionIds
	} else {
		smallQuestion := data.(models.SmallQuestion)
		tempMap := make(map[int]string)
		tempMap[smallQuestion.RealType] = models.SmallQuestionType[smallQuestion.RealType]
		if smallQuestion.RealType == models.OBJECTIVELY_BLANK || smallQuestion.RealType == models.SUBJECTIVITY_BLANK {
			tempMap[models.OBJECTIVELY_BLANK] = "客观填空题"
			tempMap[models.SUBJECTIVITY_BLANK] = "主观填空题"

			//如果这道题是客观填空题，而且还没答案
			if smallQuestion.RealType == models.OBJECTIVELY_BLANK {
				if len(smallQuestion.RealAnswer) == 0 {
					input := models.FindBlankNum(smallQuestion.Content)
					smallQuestion.RealAnswer = make([]string, input)
					data = smallQuestion
				}
			}
		}
		self.Data["QuestionTypeMap"] = tempMap
	}

	self.Data["IsBig"] = isBig
	self.Data["Data"] = data
	self.display()
}

func (self *QuestionController) Edit() {
	questionId, _ := self.GetInt64("question_id")
	isBig, _ := self.GetBool("isBig")
	questionType, _ := self.GetInt("s_question_type")
	data := make(map[string]interface{})

	if isBig {
		content := strings.TrimSpace(self.GetString("question_content"))
		data["content"] = content

	} else {
		content := strings.TrimSpace(self.GetString("question_content"))
		solution := strings.TrimSpace(self.GetString("s_question_solution"))
		score, _ := self.GetFloat("s_question_score", -100)
		difficulty, _ := self.GetFloat("s_question_difficulty", -100)

		option_len, _ := self.GetInt("op_len", -1)
		options := make(map[int]string)

		if option_len != -1 {
			for i := 0; i < option_len; i++ {
				result := strings.TrimSpace(self.GetString("options" + strconv.Itoa(i)))
				if len(result) > 0 {
					options[i] = result
				}
			}
		}

		an_len, _ := self.GetInt("an_len", -1)
		answers := make(map[int]string)

		if an_len != -1 {
			for i := 0; i < an_len; i++ {
				result := strings.TrimSpace(self.GetString("answers" + strconv.Itoa(i)))
				if len(result) > 0 {
					answers[i] = result
				}
			}
		}

		data["content"] = content
		data["solution"] = solution
		data["score"] = score
		data["difficulty"] = difficulty
		data["options"] = options
		data["answers"] = answers
		data["an_len"] = an_len
	}

	if err := models.UpdateQuestion(self.user, questionId, isBig, questionType, data); err != nil {
		self.ajaxMsg(err.Error(), -1)
	}
	self.ajaxMsg("", 0)
}
