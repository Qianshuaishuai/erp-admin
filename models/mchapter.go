package models

import (
	"errors"
)

func UpdateChapterByIndex(setId int64, index int, chapterName string, chapterDetail string, chapterQuestionCount int, chapterScore float64) error {
	tx := GetDb().Begin()
	updated := make(map[string]interface{})

	if len(chapterName) > 0 {
		updated["F_name"] = chapterName
	}

	if len(chapterDetail) > 0 {
		updated["F_detail"] = chapterDetail
	}

	if chapterQuestionCount != -100 {
		updated["F_question_count"] = chapterQuestionCount
	}

	if chapterScore != -100 {
		updated["F_preset_score"] = chapterScore
	}

	var resultChapters []PaperQuestionSetChapter

	if len(updated) == 0 {
		tx.Rollback()
		return errors.New("没有更新任何数据")
	}

	err := tx.Table("t_paper_question_set_chapters").Where("F_set_id = ?", setId).Find(&resultChapters).Error
	if err != nil {
		tx.Rollback()
		return errors.New("查询SetID" + string(setId) + " 失败:" + err.Error())
	}

	targetChapter := resultChapters[index]

	row := tx.Table("t_paper_question_set_chapters").
		Where("F_set_id = ? AND F_name = ? AND F_detail = ? AND F_preset_score = ? AND F_question_count = ?",
		setId, targetChapter.Name, targetChapter.Detail, targetChapter.PresetScore, targetChapter.QuestionCount).
		UpdateColumns(updated).RowsAffected

	if row != 1 {
		err := tx.Error
		tx.Rollback()
		if err != nil {
			return errors.New("更新失败 rows" + string(row) + " 失败:" + err.Error())
		} else {
			return errors.New("更新失败 rows" + string(row))
		}
	}

	tx.Commit()
	return nil
}
