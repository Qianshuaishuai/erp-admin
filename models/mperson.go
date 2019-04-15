package models

import (
	"errors"
	"strings"
)

func GetPersonTagListSimple(q string, limit int, page int) (list []PersonTag, count int64) {
	db := GetEliteDb().Table("t_person_tags")
	queryParams := make(map[string]interface{})
	list = make([]PersonTag, 0)

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

	db = db.Model(PersonTag{}).
		Where(queryParams).
		Count(&count)

	db.Limit(limit).
		Offset(offset).
		Scan(&list)

	return
}

func DeletePersonTag(id int) error {
	tx := GetEliteDb().Begin()
	tx = tx.Table("t_person_tags")

	count := tx.Where("id = ?", id).Delete(PersonTag{}).RowsAffected

	if count <= 0 {
		tx.Rollback()
		return errors.New("更新失败")
	}

	err := tx.Table("t_user_tags").Where("person = ?", id).Delete(UserTag{}).Error

	if err != nil {
		tx.Rollback()
		return errors.New("更新失败")
	}

	tx.Commit()

	return nil
}

func AddPersonTag(name, icon string) error {
	db := GetEliteDb().Table("t_person_tags")

	var tag PersonTag
	tag.Name = name
	tag.Icon = icon

	count := db.Create(&tag).RowsAffected

	if count <= 0 {
		return errors.New("添加失败")
	}

	return nil
}
