package models

import (
	"errors"
	"strings"
)

func GetIndustryTagListSimple(q string, limit int, page int) (list []IndustryTag, count int64) {
	db := GetEliteDb().Table("t_industry_tags")
	queryParams := make(map[string]interface{})
	list = make([]IndustryTag, 0)

	//处理分页参数
	var offset int
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
	}

	// 将搜索字符串按空格拆分
	q = strings.TrimSpace(q)
	var qstring string
	if len(q) > 0 {
		qs := strings.Fields(q)
		for _, v := range qs {
			qstring += "%" + v
		}
		qstring += "%"
	}

	if len(qstring) > 0 {
		db = db.Where("name LIKE ?", qstring)
	}

	db = db.Model(IndustryTag{}).
		Where(queryParams).
		Count(&count)

	db.Limit(limit).
		Offset(offset).
		Scan(&list)

	return
}

func DeleteIndustryTag(id int) error {
	tx := GetEliteDb().Begin()
	tx = tx.Table("t_industry_tags")

	count := tx.Where("id = ?", id).Delete(IndustryTag{}).RowsAffected

	if count <= 0 {
		tx.Rollback()
		return errors.New("更新失败")
	}

	err := tx.Table("t_user_industrys").Where("industry = ?", id).Delete(UserIndustry{}).Error

	if err != nil {
		tx.Rollback()
		return errors.New("更新失败")
	}

	err = tx.Table("t_project_industrys").Where("industry = ?", id).Delete(ProjectIndustry{}).Error

	if err != nil {
		tx.Rollback()
		return errors.New("更新失败")
	}

	tx.Commit()

	return nil
}

func AddIndustryTag(name string) error {
	db := GetEliteDb().Table("t_industry_tags")

	var tag IndustryTag
	tag.Name = name

	count := db.Create(&tag).RowsAffected

	if count <= 0 {
		return errors.New("添加失败")
	}

	return nil
}
