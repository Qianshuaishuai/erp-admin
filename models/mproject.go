package models

import (
	"errors"
	"strings"
)

func GetProjectListSimple(q string, limit int, page int, sort int, phone int64) (list []Project, count int64) {
	db := GetEliteDb().Table("t_projects")
	queryParams := make(map[string]interface{})
	list = make([]Project, 0)

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

	var sortStr = "DESC" // 默认时间 降序
	if sort == 1 {
		sortStr = "ASC"
	}

	db = db.Model(Project{}).
		Where(queryParams).
		Count(&count)

	db.Limit(limit).
		Offset(offset).
		Order("time " + sortStr).
		Scan(&list)
	return
}

func ChangeProjectStatus(id int, status int) error {
	db := GetEliteDb().Table("t_projects")

	count := db.Where("id = ?", id).Update("status", status).RowsAffected

	if count <= 0 {
		return errors.New("更新失败")
	}

	return nil
}
