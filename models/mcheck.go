package models

import (
	"time"
	"encoding/json"
	"bytes"
)

const (
	DATA_TYPE_PAPER    = 1
	DATA_TYPE_CHAPTER  = 2
	DATA_TYPE_SMALL_QUESTION = 3
	DATA_TYPE_BIG_QUESTION = 4
)

const (
	OP_EDIT   = 101
	OP_ADD    = 201
	OP_DELETE = 301
)

type HistoryDetail struct {
	FieldName string
	Old       string
	New       string
}

type CheckData struct {
	ModifyId    int64     `gorm:"primary_key;column:F_modify_id;type:BIGINT(20)"`
	ModifyDate  time.Time `gorm:"column:F_modify_date;type:DATETIME"`
	ModifyAdmin string    `gorm:"column:F_modify_admin;"`
	DeletedAt   time.Time

	DataId      int64  `gorm:"column:F_data_id;type:BIGINT(20)"`
	DataType    int    `gorm:"column:F_date_type;type:TINYINT"`
	DataOperate int    `gorm:"column:F_date_operate;type:SMALLINT"`
	CommitFlag  int    `gorm:"column:F_commit_flag;type:TINYINT"`
	DetailsDB   string `gorm:"column:F_detail_db;type:TEXT"`
}

func AddOperateData(resId int64, resType, resOperate int, details []HistoryDetail) {
	tx := GetDb().Begin()
	var checkData CheckData
	var snowCurl MSnowflakCurl
	//生成ModifyId
	checkData.ModifyId = int64(snowCurl.GetIntId(true))
	checkData.ModifyDate = time.Now()
	checkData.ModifyAdmin = "Admin"

	checkData.DataId = resId
	checkData.DataType = resType
	checkData.DataOperate = resOperate

	checkData.CommitFlag = 0
	var buffer bytes.Buffer
	jEncoder := json.NewEncoder(&buffer)
	jEncoder.SetEscapeHTML(false)
	jEncoder.Encode(details)
	checkData.DetailsDB = buffer.String()

	err := tx.Create(&checkData).Error

	if err != nil {
		GetLogger().LogErr(err, "add_operate_data")
		tx.Rollback()
		return
	}
	tx.Commit()
}

func GetCheckList(limit int, page int) (result []CheckData, count int64) {
	result = make([]CheckData, 0)

	//处理分页参数
	var offset int
	if limit > 0 && page > 0 {
		offset = (page - 1) * limit
	}

	GetDb().Table("t_check_data").Count(&count).Limit(limit).Offset(offset).Order("F_modify_date DESC").Scan(&result)

	return
}
