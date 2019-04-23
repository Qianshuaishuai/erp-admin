package models

import (
	"elite-admin/helper"
	"errors"
	"strings"
	"time"
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

func GetProjectDetail(id int64) (data ProjectDetail, err error) {
	GetEliteDb().Table("t_projects").Where("id = ?", id).Find(&data.Project)

	var industryIDs []int
	GetEliteDb().Table("t_project_industrys").Where("project = ?", id).Pluck("industry", &industryIDs)

	GetEliteDb().Table("t_industry_tags").Where("id in (?)", industryIDs).Pluck("name", &data.Industrys)

	return data, nil
}

func DeleteProjectDetail(id int) (err error) {
	tx := GetEliteDb().Begin()
	err = tx.Table("t_projects").Where("id = ?", id).Delete(Project{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Table("t_project_industrys").Where("project = ?", id).Delete(ProjectIndustry{}).Error

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func AddProject(name, typeName, address, money, agency, introduce, addtip, idcard string, phone int, cardImageURL, backgroundImageURL string, industrys string) error {
	var snowCurl MSnowflakCurl
	var project Project
	project.Name = name
	project.Type = typeName
	project.Address = address
	project.Money = money
	project.Agency = agency
	project.Introduce = introduce
	project.AddTip = addtip
	project.IDCard = idcard
	project.Phone = phone
	project.Icon = cardImageURL
	project.Background = backgroundImageURL
	project.Time = time.Now()
	project.Status = 1
	project.ID = snowCurl.GetIntId(true)

	err := GetEliteDb().Table("t_projects").Create(&project).Error

	if err != nil {
		return errors.New("添加失败")
	}

	industryIDs, _ := helper.TransformStringToInt64Arr(industrys)
	industryTags, _ := GetIndustryTagListSimple("", 100, 1, 0)

	for i := range industryIDs {
		var industry ProjectIndustry
		industry.Industry = industryTags[industryIDs[i]-1].ID
		industry.Project = project.ID
		GetEliteDb().Table("t_project_industrys").Create(&industry)
	}

	return nil
}
