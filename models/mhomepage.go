package models

import (
	"errors"
	"strings"
)

func GetHomePageTagListSimple(q string, limit int, page int) (list []HomeShowDetail, count int64) {
	db := GetEliteDb().Table("t_home_shows")
	queryParams := make(map[string]interface{})
	list = make([]HomeShowDetail, 0)
	mlist := make([]HomeShow, 0)

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

	db = db.Model(HomeShow{}).
		Where(queryParams).
		Count(&count)

	db.Limit(limit).
		Offset(offset).
		Scan(&mlist)

	var contents []HomeContent
	GetEliteDb().Table("t_home_contents").Find(&contents)

	for m := range mlist {
		var detail HomeShowDetail
		detail.HomeShow = mlist[m]

		for c := range contents {
			if contents[c].ID == detail.HomeShow.ID {
				detail.HomeContents = append(detail.HomeContents, contents[c])
			}
		}

		list = append(list, detail)
	}
	return
}

func DeleteHomeShow(id int) error {
	tx := GetEliteDb().Begin()
	tx = tx.Table("t_home_shows")

	count := tx.Where("id = ?", id).Delete(HomeShow{}).RowsAffected

	if count <= 0 {
		tx.Rollback()
		return errors.New("更新失败")
	}

	err := tx.Table("t_home_contents").Where("id = ?", id).Delete(HomeContent{}).Error

	if err != nil {
		tx.Rollback()
		return errors.New("更新失败")
	}

	tx.Commit()

	return nil
}

func GetHomePageType() (datas []HomePageType) {
	var typeSimpleA HomePageType
	typeSimpleA.ID = 1
	typeSimpleA.Name = "正文标签"

	var typeSimpleB HomePageType
	typeSimpleB.ID = 2
	typeSimpleB.Name = "外链标签"

	datas = append(datas, typeSimpleA)
	datas = append(datas, typeSimpleB)

	return datas

}

func AddNewHomeShow(name string, typeID int, iconImageURL string, URL string, datas []HomePageContent) error {
	var homeShow HomeShow
	homeShow.Name = name
	homeShow.Icon = iconImageURL
	homeShow.Type = typeID
	if typeID == 2 {
		homeShow.URL = URL
	}

	err := GetEliteDb().Table("t_home_shows").Create(&homeShow).Error
	if err != nil {
		return errors.New("创建失败")
	}

	if typeID == 1 {
		for d := range datas {
			var content HomeContent
			content.ID = homeShow.ID
			content.Content = datas[d].Content
			content.Type = datas[d].TypeID
			content.Index = d + 1
			GetEliteDb().Table("t_home_contents").Create(&content)
		}
	}

	return nil
}
