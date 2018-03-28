package models

import (
	"time"
	"strings"
	"errors"
	"strconv"
	"dreamEbagPaperAdmin/helper"
)

type PaperType struct {
	Id   int8   `gorm:"primary_key;column:F_paper_type_id;type:TINYINT(4)" json:"id"`
	Name string `gorm:"column:F_name;size:20" json:"name"`
}

type Province struct {
	ProvinceId uint   `gorm:"column:F_province_id;primary_key;" json:"id"`
	Name       string `gorm:"column:F_name;size:6;unique" json:"name"`
}

//试卷表的详细信息 （从试卷表中拆分出来）
type PaperQuestionSet struct {
	SetId            int64  `gorm:"primary_key;column:F_set_id;type:BIGINT(20)" json:"id"`
	PaperId          int64  `gorm:"column:F_paper_id;type:BIGINT(20)" json:"paperId"`
	PaperName        string `gorm:"column:F_paper_name;size:100" json:"paperName"`
	TimeToAccomplish uint   `gorm:"column:F_time_to_accomplish" json:"timeToAccomplish"` //试卷完成时间

	QuestionIds              string                    `gorm:"column:F_question_ids;type:TEXT" json:"questionIds"` // 问题ID列表 字符串 按试题在试卷中顺序排列
	PaperQuestionSetChapters []PaperQuestionSetChapter `gorm:"ForeignKey:SetId;" json:"paperQuestionSetChapters"`  //试卷包含的章节（详细）
}

//试卷章节介绍
type PaperQuestionSetChapter struct {
	Name             string  `gorm:"column:F_name;size:45" json:"name"`
	Detail           string  `gorm:"column:F_detail;type:TEXT" json:"desc" mapstructure:"desc"` //本章说明/介绍
	QuestionCount    uint    `gorm:"column:F_question_count;" json:"questionCount"`             //包含的问题个数（可能大于实际题目数量）
	TimeToAccomplish uint    `gorm:"column:F_time;" json:"time"`                                //完成所需的时间
	PresetScore      float32 `gorm:"column:F_preset_score;" json:"presetScore"`                 //分数(有些章节题目有缺失，请注意计算时使用题目实际数量)

	SetId            int64         `gorm:"column:F_set_id;type:BIGINT(20)" json:"setId"`
	QuestionsContent []interface{} `gorm:"-" json:"questionsContent"`
}

type Paper struct {
	PaperId    int64     `gorm:"primary_key;column:F_paper_id;type:BIGINT(20)" json:"id"`
	Name       string    `gorm:"column:F_name;size:80" json:"name"`
	PaperType  int8      `gorm:"column:F_paper_type;type:TINYINT(4)" json:"type" mapstructure:"type"` //试卷类型 真题 or 模拟题
	Difficulty float32   `gorm:"column:F_difficulty" json:"difficulty"`                               // 难度范围2-7
	FullScore  uint      `gorm:"column:F_full_score" json:"fullScore"`                                // 总分
	Date       time.Time `gorm:"column:F_date;type:DATE" json:"date"`                                 //试卷的编写日期

	SemesterId  uint             `gorm:"column:F_semester_id" json:"semesterId"`              //试卷对应的学期
	CourseId    uint             `gorm:"column:F_course_id;type:TINYINT(2);" json:"courseId"` //试卷对应的课程
	Provinces   []Province       `gorm:"many2many:t_paper_province;" json:"provinces"`        // 试卷适用的省份
	QuestionSet PaperQuestionSet `gorm:"ForeignKey:PaperId;" json:"questionSet"`
}

type PaperSimple struct {
	PaperId       int64     `gorm:"primary_key;column:F_paper_id;type:BIGINT(20)" json:"id"`
	Name          string    `gorm:"column:F_name;size:80" json:"name"`
	PaperType     int8      `gorm:"column:F_paper_type;type:TINYINT(4)" json:"type"` //试卷类型 真题 or 模拟题
	PaperTypeName string    `gorm:"-" json:"typeName"`                               //试卷类型的名称（描述）
	Date          time.Time `gorm:"column:F_date;type:DATE" json:"date"`             //试卷的编写日期
}

func GetPaperListSimple(q string, limit int, page int, sort int) (list []PaperSimple, count int64) {
	db := GetDb().Table("t_papers")
	queryParams := make(map[string]interface{})
	list = make([]PaperSimple, 0)

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
		db = db.Where("F_name LIKE ?", qstring)
	}

	var sortStr = "DESC" // 默认时间 降序
	if sort == 1 {
		sortStr = "ASC"
	}

	db = db.Select("F_paper_id,F_name,F_paper_type,F_date").
		Where(queryParams).
		Count(&count)

	db.Limit(limit).
		Offset(offset).
		Order("F_date " + sortStr).
		Scan(&list)

	//根据所查到的试卷的type找到对应的PaperType
	findPaperType(list)
	return
}

func findPaperType(list []PaperSimple) {
	for i := range list {
		paperType := GetPaperType(list[i].PaperType)
		if paperType.Id == 0 {
			//没找到
			paperType.Name = ""
		}
		list[i].PaperTypeName = paperType.Name
	}
}

func GetProvinces() (provinces []Province) {
	GetDb().Find(&provinces)
	return provinces
}

func GetPaperType(typeId int8) PaperType {
	var paperType PaperType
	GetDb().Find(&paperType, typeId)
	return paperType
}

func GetAllPaperType() []PaperType {
	var paperTypes []PaperType
	GetDb().Find(&paperTypes)
	return paperTypes
}

//根据resourceId去t_papers表找
func GetPaper(resourceId int64) (Paper) {
	var info Paper
	var isRecordNotFound bool

	if resourceId > 0 {
		isRecordNotFound = GetDb().Preload("Provinces").Find(&info, resourceId).RecordNotFound()

		if isRecordNotFound {
			//没找到试卷

		} else {
			var paperQuestionSet PaperQuestionSet
			GetDb().Preload("PaperQuestionSetChapters").Find(&paperQuestionSet, "F_paper_id = ?", resourceId)

			info.QuestionSet = paperQuestionSet
		}
	}
	return info
}

func UpdatePaper(paperId int64, paperName string, fullScore int, paperType int, difficulty float64, provinces string) error {
	//先处理省份
	tx := GetDb().Begin()

	if len(provinces) > 0 {
		provincesSplit := strings.Split(provinces, ",")
		if len(provincesSplit) > 0 {
			//删除这个试卷的关联省份
			err := tx.Exec("DELETE FROM t_paper_province WHERE paper_F_paper_id = ?", paperId).Error

			if err != nil {
				tx.Rollback()
				return errors.New("删除省份失败:" + err.Error())
			}

			for i := range provincesSplit {
				if len(provincesSplit[i]) > 0 {
					newProvinceId, _ := strconv.ParseInt(provincesSplit[i], 10, 64)
					if newProvinceId != 0 {
						err = tx.Exec("INSERT INTO t_paper_province VALUES (?,?)", paperId, newProvinceId).Error
						if err != nil {
							tx.Rollback()
							return errors.New("插入省份失败:" + err.Error())
						}
					}
				}
			}
		}
	}

	updated := make(map[string]interface{})

	//处理paperName
	if len(paperName) > 0 {
		updated["F_name"] = paperName
	}

	if fullScore != -100 {
		updated["F_full_score"] = fullScore
	}

	if paperType != -100 {
		updated["F_paper_type"] = paperType
	}

	if difficulty != -100 {
		updated["F_difficulty"] = difficulty
	}

	//更新时间
	updated["F_date"] = time.Now()

	err := tx.Model(&Paper{}).Where("F_paper_id = ?", paperId).Updates(updated).Error
	if err != nil {
		tx.Rollback()
		return errors.New("更新试卷失败:" + err.Error())
	}

	tx.Commit()
	return nil
}

//找到resId中chapterQuestionCount和q指向的部分
func GetTheQuestionByQ(resIds []int64, q int, chapterQuestionCount int) (startIndex, endIndex int) {
	startIndex = q
	for {
		//先瞅瞅q当前指向的题目有没有
		if resIds[q] == 0 {
			chapterQuestionCount--
		} else {
			//这题是大题是小题
			if isBig, bigCount := GetQuestionTranslateTypeById(resIds[q]); isBig {
				chapterQuestionCount -= bigCount
			} else {
				chapterQuestionCount--
			}
		}
		q++
		if chapterQuestionCount == 0 {
			break
		}
	}
	endIndex = q
	return
}

//根据id确定这道题目是大题还是小题
func GetQuestionTranslateTypeById(resId int64) (isBig bool, bigCount int) {
	var s []string
	GetDb().Table("t_large_questions").Where("F_big_question_id = ?", resId).Pluck("F_question_ids", &s)

	if len(s) > 0 {
		resIds, _ := helper.TransformStringToInt64Arr(s[0])
		isBig = true
		bigCount = len(resIds)
	} else {
		isBig = false
		bigCount = 0
	}
	return
}
