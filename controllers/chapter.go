package controllers

import (
	"strings"
	"dreamEbagPaperAdmin/models"
)

type ChapterController struct {
	BaseController
}

func (self *ChapterController) Edit() {
	Set_id, _ := self.GetInt64("set_id")
	Index, _ := self.GetInt("index")
	if Set_id != 0 {
		chapter_name := strings.TrimSpace(self.GetString("chapter_name"))
		chapter_detail := strings.TrimSpace(self.GetString("chapter_detail"))
		chapter_question_count, _ := self.GetInt("chapter_question_count", -100)
		chapter_score, _ := self.GetFloat("chapter_score", -100.0)

		err := models.UpdateChapterByIndex(Set_id, Index, chapter_name, chapter_detail, chapter_question_count, chapter_score)

		if err != nil {
			self.ajaxMsg(err.Error(), -1)
		}

		self.ajaxMsg("", 0)
	}
}

