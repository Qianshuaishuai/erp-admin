package models

import (
	"errors"
	"strconv"
	"strings"
)

func GetConnectionListSimple(q string, limit int, page int, sort int, phone int64) (list []Connection, count int64) {
	db := GetEliteDb().Table("t_connections")
	queryParams := make(map[string]interface{})
	list = make([]Connection, 0)

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

	// if len(qstring) > 0 {
	// 	db = db.Where("username LIKE ?", qstring)
	// }

	var sortStr = "DESC" // 默认时间 降序
	if sort == 1 {
		sortStr = "ASC"
	}

	db = db.Model(Connection{}).
		Where(queryParams).
		Count(&count)

	db.Limit(limit).
		Offset(offset).
		Order("time " + sortStr).
		Scan(&list)
	return
}

func ChangeConnectionStatus(phone int, status int) error {
	db := GetEliteDb().Table("t_connections")

	count := db.Where("phone = ?", phone).Update("status", status).RowsAffected

	if count <= 0 {
		return errors.New("更新失败")
	}

	return nil
}
