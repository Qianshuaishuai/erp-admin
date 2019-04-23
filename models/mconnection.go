package models

import (
	"elite-admin/helper"
	"errors"
	"strconv"
	"strings"
	"time"
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

func DeleteConnectionSimple(phone int) (err error) {
	tx := GetEliteDb().Begin()
	err = tx.Table("t_connections").Where("phone = ?", phone).Delete(Connection{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Table("t_users").Where("phone = ?", phone).Delete(UserInfoSimple{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Table("t_user_tags").Where("phone = ?", phone).Delete(UserTag{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func GetConnectionDetail(phone int64) (data UserInfoDetail, err error) {
	GetEliteDb().Table("t_connections").Where("phone = ?", phone).Find(&data.Connection)
	GetEliteDb().Table("t_users").Where("phone = ?", phone).Find(&data.UserInfoSimple)

	var tagIDs []int
	GetEliteDb().Table("t_user_tags").Where("phone = ?", phone).Pluck("person", &tagIDs)

	var tagName []string
	GetEliteDb().Table("t_person_tags").Where("id in (?)", tagIDs).Pluck("name", &tagName)
	data.Tags = tagName

	return data, nil
}

func AddConnection(phone, look, good int, username, job, position, profess, agency, address, introduce, achieve, school, iconImageURL, cardImageURL string, tags string) error {
	tx := GetEliteDb().Begin()

	var userInfo UserInfoSimple
	userInfo.Phone = phone
	userInfo.Password = "111111"
	userInfo.Username = username
	userInfo.Job = job
	userInfo.Position = position
	userInfo.Profess = profess
	userInfo.Agency = agency
	userInfo.Address = address
	userInfo.Introduce = introduce
	userInfo.Achieve = achieve
	userInfo.School = school
	userInfo.Icon = iconImageURL
	userInfo.Register = time.Now()

	var count int
	tx.Table("t_users").Where("phone = ?", phone).Count(&count)

	if count > 0 {
		tx.Rollback()
		return errors.New("手机号已存在")
	}

	err := tx.Table("t_users").Create(&userInfo).Error

	if err != nil {
		tx.Rollback()
		return errors.New("创建失败")
	}

	var connection Connection
	connection.Card = cardImageURL
	connection.Good = good
	connection.Look = look
	connection.Phone = phone
	connection.Status = 1
	connection.Time = time.Now()
	connection.From = 1

	var aCount int
	tx.Table("t_connections").Where("phone = ?", phone).Count(&aCount)

	if aCount > 0 {
		tx.Rollback()
		return errors.New("手机号已存在")
	}

	err = tx.Table("t_connections").Create(&connection).Error

	if err != nil {
		tx.Rollback()
		return errors.New("创建失败")
	}

	tagIDs, _ := helper.TransformStringToInt64Arr(tags)
	personTags, _ := GetPersonTagListSimple("", 100, 1)

	for t := range tagIDs {
		var userTag UserTag
		userTag.Person = personTags[tagIDs[t]-1].ID
		userTag.Phone = phone
		GetEliteDb().Table("t_user_tags").Create(&userTag)
	}

	tx.Commit()
	return nil
}
