package models

import (
	"strings"
	"errors"
)

func SaveAddPaperTemp(name string, fullScore int, timeToAccomplish int, paperYear int,
	courseId int, semesterId int, typeId int, difficulty float64, provinceIds string) error {
	var snowCurl MSnowflakCurl
	add := AddPaperTemp{
		PaperId:          int64(snowCurl.GetIntId(true)),
		Name:             name,
		FullScore:        fullScore,
		TimeToAccomplish: timeToAccomplish,
		PaperYear:        paperYear,
		CourseId:         courseId,
		SemesterId:       semesterId,
		TypeId:           typeId,
		Difficulty:       difficulty,
		ProvinceIds:      provinceIds,
		Status:           ADDPAPERTEMP_STATUS_EDIT,
	}

	err := GetDb().Create(&add).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAddPaperTempList(limit int, page int, cFlag, uFlag bool, cSort, uSort int) (list []AddPaperTemp, count int64) {
	db := GetDb().Model(&AddPaperTemp{})
	list = make([]AddPaperTemp, 0)

	//处理分页参数
	var offset int
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
	}

	db = db.
		Count(&count).
		Limit(limit).
		Offset(offset)

	// ASC DESC 排序字段拼接
	if cFlag {
		var sortString = "DESC"
		if cSort == 1 {
			sortString = "ASC"
		}
		db = db.Order("created_at " + sortString)
	}

	if uFlag {
		var sortString = "DESC"
		if uSort == 1 {
			sortString = "ASC"
		}
		db = db.Order("updated_at " + sortString)
	}
	db.Scan(&list)
	return
}

func DeleteAddPaper(paperId int64) error {
	if paperId != 0 {
		db := GetDb()
		return db.Delete(&AddPaperTemp{}, "F_paper_id = ?", paperId).Error
	}
	return errors.New("PaperId 为 0")
}

func GetAddPaperTemp(paperId int64) (info AddPaperTemp) {
	if paperId > 0 {
		GetDb().Find(&info, paperId)
	}
	return
}

func UpdateAddPaper(
	paperId int64,
	paperName string,
	fullScore int,
	paperTime int,
	paperCourse int,
	paperSemester int,
	paperType int,
	difficulty float64,
	provinces string,
) error {
	tx := GetDb().Begin()
	//先处理省份
	provinces = strings.TrimRight(provinces, ",")

	updated := make(map[string]interface{})
	//处理paperName
	if len(paperName) > 0 {
		updated["F_name"] = paperName
	}

	if fullScore != -100 {
		updated["F_full_score"] = fullScore
	}

	if paperTime != -100 {
		updated["F_time"] = paperTime
	}

	if paperCourse != -100 {
		updated["F_course_id"] = paperCourse
	}

	if paperSemester != -100 {
		updated["F_semester_id"] = paperSemester
	}

	if paperType != -100 {
		updated["F_type_id"] = paperType
	}

	if difficulty != -100 {
		updated["F_difficulty"] = difficulty
	}

	//更新省份
	updated["F_province_ids"] = provinces

	err := tx.Model(&AddPaperTemp{}).Where("F_paper_id = ?", paperId).Updates(updated).Error
	if err != nil {
		return HandleErrByTx(errors.New("更新试卷失败:"+err.Error()), tx)
	}

	tx.Commit()
	return nil
}
