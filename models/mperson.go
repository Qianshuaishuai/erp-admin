package models

import (
	"errors"
	"strings"
)

func GetPersonTagListSimple(q string, limit int, page int, sort int) (list []PersonTag, count int64) {
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

	var sortStr = "DESC" // 默认时间 降序
	if sort == 1 {
		sortStr = "ASC"
	}

	if len(qstring) > 0 {
		db = db.Where("name LIKE ?", qstring)
	}

	db = db.Model(PersonTag{}).
		Where(queryParams).
		Count(&count)

	db.Limit(limit).
		Offset(offset).
		Order("plain " + sortStr).
		Scan(&list)

	return
}

func GetPersonDetail(id int64) (data PersonTag, err error) {
	GetEliteDb().Table("t_person_tags").Where("id = ?", id).Find(&data)

	return data, nil
}

func EditPersonDetail(id, index int, newName, imageURL string) (err error) {
	var tag PersonTag
	tag.ID = int(id)
	tag.Name = newName
	tag.Plain = index
	tag.Icon = imageURL

	GetEliteDb().Table("t_person_tags").Where("id = ?", id).Update(&tag)

	return nil
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
