package models

import (
	"strconv"
	"strings"
)

func GetUserListSimple(q string, limit int, page int, sort int, phone int64) (list []UserInfoSimple, count int64) {
	db := GetEliteDb().Table("t_users")
	queryParams := make(map[string]interface{})
	list = make([]UserInfoSimple, 0)

	//处理分页参数
	var offset int
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
	}

	//查手机号
	if phone > 0 {
		db = db.Where("phone like ?", "%"+strconv.FormatInt(phone, 10)+"%")
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
		db = db.Where("username LIKE ?", qstring)
	}

	var sortStr = "DESC" // 默认时间 降序
	if sort == 1 {
		sortStr = "ASC"
	}

	db = db.Model(UserInfoSimple{}).
		Where(queryParams).
		Count(&count)

	db.Limit(limit).
		Offset(offset).
		Order("register " + sortStr).
		Scan(&list)

	return
}
