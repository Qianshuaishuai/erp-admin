package models

import (
	"dreamEbagPaperAdmin/helper"
	"encoding/json"
	"strings"
	"errors"
	"bytes"
	"regexp"
	"strconv"
)

const (
	RADIO_CHOICE         = 10001
	MULIT_CHOICE         = 10002
	INDETERMINATE_CHOICE = 10003
	JUDGE_CHOICE         = 10004
	OBJECTIVELY_BLANK    = 10005
	SUBJECTIVITY_BLANK   = 10006
	QA_BLANK             = 10007
)

var SmallQuestionType map[int]string

func init() {
	SmallQuestionType = make(map[int]string)

	SmallQuestionType[RADIO_CHOICE] = "单选题"
	SmallQuestionType[MULIT_CHOICE] = "多选题"
	SmallQuestionType[INDETERMINATE_CHOICE] = "不定项选择题"
	SmallQuestionType[JUDGE_CHOICE] = "判断题"
	SmallQuestionType[OBJECTIVELY_BLANK] = "客观填空题"
	SmallQuestionType[SUBJECTIVITY_BLANK] = "主观填空题"
	SmallQuestionType[QA_BLANK] = "问答题"
}

// 映射小题问题类型
func MapTypeSmallQuestion(oldType int) (newType int) {
	switch oldType {
	case 1, 4, 6:
		newType = RADIO_CHOICE
	case 2:
		newType = MULIT_CHOICE
	case 3:
		newType = INDETERMINATE_CHOICE
	case 5:
		newType = JUDGE_CHOICE
	case 61, 64:
		newType = OBJECTIVELY_BLANK
	case 16, 63, 50:
		newType = SUBJECTIVITY_BLANK
	case 13, 14, 15, 65:
		newType = QA_BLANK
	default:
		newType = 0
	}
	return
}

// 反向映射小题问题类型
func MapTypeSmallQuestionRevert(newType int) (oldType int) {
	switch newType {
	case RADIO_CHOICE:
		return 1
	case MULIT_CHOICE:
		return 2
	case INDETERMINATE_CHOICE:
		return 3
	case JUDGE_CHOICE:
		return 5
	case OBJECTIVELY_BLANK:
		return 61
	case SUBJECTIVITY_BLANK:
		return 60
	case QA_BLANK:
		return 13
	}
	return 0
}

type BigQuestion struct {
	QuestionId     int64  `gorm:"primary_key;column:F_big_question_id;type:BIGINT(20)" json:"id"`
	Content        string `gorm:"column:F_content;type:LONGTEXT" json:"content"` //问题的内容
	BigQuestionIds string `gorm:"column:F_question_ids;type:TEXT" json:"-"`

	RealType   int
	Options    []string
	RealAnswer []string
	Solution   string  `gorm:"column:F_solution;type:LONGTEXT" json:"solution"`
	Score      float32 `gorm:"column:F_score;type:FLOAT(4,1)" json:"score"`
	Difficulty float64 `gorm:"column:F_difficulty;" json:"difficulty"`
}

//问题的知识点
type Keypoint struct {
	KeypointId int64  `gorm:"primary_key;column:F_keypoint_id;type:BIGINT(20)" json:"id"`
	Name       string `gorm:"column:F_name;size:255" json:"name"`
	Type       int    `gorm:"column:F_type;" json:"type"`
}

type SmallQuestion struct {
	QuestionId    int64   `gorm:"primary_key;column:F_question_id;type:BIGINT(20)" json:"id"`
	Content       string  `gorm:"column:F_content;type:LONGTEXT" json:"content"`     //问题的内容
	Score         float32 `gorm:"column:F_score;type:FLOAT(4,1)" json:"score"`       //问题分数 最大值999.9
	Accessories   string  `gorm:"column:F_accessories;type:TEXT" json:"accessories"` //问题的附加内容 （选择题选项等）
	Solution      string  `gorm:"column:F_solution;type:LONGTEXT" json:"solution"`   // 问题的解答
	Difficulty    float64 `gorm:"column:F_difficulty;" json:"difficulty"`
	CorrectAnswer string  `gorm:"column:F_correct_answer;type:TEXT" json:"correctAnswer"` //正确答案 （不一定有）
	Type          int     `gorm:"column:F_type" json:"type"`
	RealType      int
	Options       []string
	RealAnswer    []string
	//Keypoints     []Keypoint `gorm:"many2many:t_keypoint_question;" json:"keypoints"`        //考察的知识点
}

func GetQuestion(resId int64, q string) (isBig bool, data interface{}) {
	isBig, _ = GetQuestionTranslateTypeById(resId)
	db := GetDb()

	if resId != 0 {
		if isBig {
			// 大题
			var bigQuestion BigQuestion
			db.Table("t_large_questions").Where("F_big_question_id = ?", resId).Scan(&bigQuestion)
			data = bigQuestion
		} else {
			//小题
			var smallQuestion SmallQuestion
			db.Table("t_questions").Where("F_question_id = ?", resId).Scan(&smallQuestion)
			smallQuestion.RealType = MapTypeSmallQuestion(smallQuestion.Type)
			smallQuestion.Options = generateOptions(smallQuestion.Accessories, smallQuestion.RealType)
			smallQuestion.RealAnswer = generateAnswer(smallQuestion.CorrectAnswer, smallQuestion.RealType)
			data = smallQuestion
		}
	} else {
		if len(q) != 0 {
			//到大题表中搜索
			q = "%" + q + "%"
			var bigQuestion BigQuestion
			db.Table("t_large_questions").Where("F_content LIKE ?", q).First(&bigQuestion)
			if bigQuestion.QuestionId != 0 {
				isBig = true
				data = bigQuestion
			} else {
				//小题表中搜索
				var smallQuestion SmallQuestion
				db.Table("t_questions").Where("F_content LIKE ?", q).Scan(&smallQuestion)
				if smallQuestion.QuestionId != 0 {
					smallQuestion.RealType = MapTypeSmallQuestion(smallQuestion.Type)
					smallQuestion.Options = generateOptions(smallQuestion.Accessories, smallQuestion.RealType)
					smallQuestion.RealAnswer = generateAnswer(smallQuestion.CorrectAnswer, smallQuestion.RealType)
				}
				isBig = false
				data = smallQuestion
			}
		}
	}

	return
}

func generateAnswer(s string, questionType int) []string {
	result := make([]string, 0)
	if len(s) > 0 {
		switch questionType {
		case RADIO_CHOICE, JUDGE_CHOICE, SUBJECTIVITY_BLANK, QA_BLANK:
			result = append(result, s)
		case MULIT_CHOICE, INDETERMINATE_CHOICE:
			res := strings.Split(s, "-")
			if len(res) > 0 {
				for i := range res {
					result = append(result, res[i])
				}
			}
		case OBJECTIVELY_BLANK:
			res := strings.Split(s, "-$-")
			if len(res) > 0 {
				for i := range res {
					result = append(result, res[i])
				}
			}
		}
	}
	return result
}

func generateOptions(s string, questionType int) (result []string) {
	result = make([]string, 0)

	if questionType == JUDGE_CHOICE {
		result = append(result, "T", "F")
		return result
	}

	if len(s) == 0 {
		return result
	}

	var temp interface{}
	jDecoder := json.NewDecoder(strings.NewReader(s))
	jDecoder.Decode(&temp)

	if dataTemp, ok := temp.(map[string]interface{}); ok {
		if e, o := dataTemp["option"]; o {
			if d, o2 := e.(map[string]interface{}); o2 {
				if x, o3 := d["options"].([]interface{}); o3 {
					//记录当前的选项
					var nowOp = 0
					for i := range x {
						xs := x[i].(string)
						if xs == "$" {
							xs := Map123toABC(nowOp)
							result = append(result, xs)
						} else {
							result = append(result, handleContent(xs))
						}
						nowOp++
					}
				}
			}
		}
	}
	return result
}

func handleContent(s string) string {
	//<tex > </tex>标签的内容中的\换成~@
	regTex, _ := regexp.Compile(`\[tex.*?](.|\n|\f|\r)*?\[\\*/tex]`)
	regla, _ := regexp.Compile(`\\+`)
	a := regTex.FindAllString(s, -1)

	for i := range a {
		x := a[i]
		x = strings.Replace(x, "\r", "\\r", -1)
		x = strings.Replace(x, "\f", "\\f", -1)
		x = strings.Replace(x, "\n", "\\n", -1)

		n := regla.ReplaceAllString(x, "~@")
		s = strings.Replace(s, a[i], n, 1)
	}
	return s
}

func handleContentRevert(s string) string {
	//<tex > </tex>标签的内容中的\换成~@
	regTex, _ := regexp.Compile(`\[tex.*?](.|\n|\f|\r)*?\[\\*/tex]`)
	regla, _ := regexp.Compile(`(~@)+`)
	a := regTex.FindAllString(s, -1)

	for i := range a {
		x := a[i]
		n := regla.ReplaceAllString(x, "\\\\")
		s = strings.Replace(s, a[i], n, 1)
	}
	return s
}

func Map123toABC(i int) string {
	var ru rune
	ru = rune(65 + i)
	return "答案" + string(ru)
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

func UpdateQuestion(questionId int64, isBig bool, questionType int, data map[string]interface{}) error {
	tx := GetDb().Begin()

	if isBig {
		//改大题表
		updated := make(map[string]interface{})
		content := data["content"].(string)

		if len(content) > 0 {
			updated["F_content"] = content
		}

		if len(updated) > 0 {
			detail := makeHistoryDetailBigQuestion(questionId, updated)

			err := tx.Table("t_large_questions").Where("F_big_question_id = ?", questionId).Updates(updated).Error
			if err != nil {
				tx.Rollback()
				return errors.New("更新大题失败:" + err.Error())
			}
			AddOperateData(questionId, DATA_TYPE_BIG_QUESTION, OP_EDIT, detail)
			tx.Commit()
		}
	} else {
		//改小题表
		updated := make(map[string]interface{})

		content := data["content"].(string)

		if len(content) > 0 {
			updated["F_content"] = content
		}

		solution := data["solution"].(string)

		if len(solution) > 0 {
			updated["F_solution"] = solution
		}

		score := data["score"].(float64)

		if score != -100 {
			updated["F_score"] = score
		}

		difficulty := data["difficulty"].(float64)

		if difficulty != -100 {
			updated["F_difficulty"] = difficulty
		}

		//处理 Options
		options := data["options"].(map[int]string)

		if len(options) > 0 {
			updated["F_accessories"] = handleContentRevert(makeNewOptions(questionId, options))
		}

		//处理Answer
		answers := data["answers"].(map[int]string)
		anLen := data["an_len"].(int)

		if len(answers) > 0 {
			updated["F_correct_answer"] = makeNewAnswers(questionId, answers, questionType, anLen)
		}

		if len(updated) > 0 {
			detail := makeHistoryDetailSmallQuestion(questionId, updated)

			err := tx.Table("t_questions").Where("F_question_id = ?", questionId).Updates(updated).Error
			if err != nil {
				tx.Rollback()
				return errors.New("更新小题失败:" + err.Error())
			}

			AddOperateData(questionId, DATA_TYPE_SMALL_QUESTION, OP_EDIT, detail)
			tx.Commit()
		}
	}
	return nil
}

func makeHistoryDetailSmallQuestion(questionId int64, updated map[string]interface{}) []HistoryDetail {
	result := make([]HistoryDetail, 0)

	for k, v := range updated {
		var temp HistoryDetail
		switch k {
		case "F_content","F_solution","F_accessories","F_correct_answer":
			temp.FieldName = k
			var z []string
			GetDb().Table("t_large_questions").Where("F_big_question_id = ?", questionId).Pluck(k, &z)
			if len(z) > 0 {
				temp.Old = z[0]
				temp.New = v.(string)
			}
		case "F_score", "F_difficulty":
			temp.FieldName = k
			var z []float64
			GetDb().Table("t_large_questions").Where("F_big_question_id = ?", questionId).Pluck(k, &z)
			if len(z) > 0 {
				temp.Old = strconv.FormatFloat(z[0], 'f', 1, 64)
				temp.New = strconv.FormatFloat(v.(float64), 'f', 1, 64)
			}
		}

		if len(temp.FieldName) > 0 {
			result = append(result, temp)
		}
	}
	return result
}

func makeHistoryDetailBigQuestion(questionId int64, updated map[string]interface{}) []HistoryDetail {
	result := make([]HistoryDetail, 0)

	for k, v := range updated {
		var temp HistoryDetail
		switch k {
		case "F_content":
			temp.FieldName = k
			var z []string
			GetDb().Table("t_large_questions").Where("F_big_question_id = ?", questionId).Pluck(k, &z)
			if len(z) > 0 {
				temp.Old = z[0]
				temp.New = v.(string)
			}
		}

		if len(temp.FieldName) > 0 {
			result = append(result, temp)
		}
	}
	return result
}

func makeNewAnswers(questionId int64, answers map[int]string, questionType int, anLen int) string {
	//查老答案
	var answersDB []string
	GetDb().Table("t_questions").Where("F_question_id = ?", questionId).Pluck("F_correct_answer", &answersDB)
	if len(answersDB) > 0 {
		ans := answersDB[0]

		switch questionType {
		case RADIO_CHOICE, JUDGE_CHOICE, QA_BLANK, SUBJECTIVITY_BLANK:
			return answers[0]
		case MULIT_CHOICE, INDETERMINATE_CHOICE:
			res := strings.Split(ans, "-")
			result := make([]string, 0, len(res))
			if len(res) > 0 {
				for i := range res {
					result = append(result, res[i])
				}
			}

			for k, v := range answers {
				result[k] = v
			}

			return strings.Join(result, "-")
		case OBJECTIVELY_BLANK:
			res := strings.Split(ans, "-$-")
			result := make([]string, 0, len(res))
			if len(res) > 0 {
				for i := range res {
					result = append(result, res[i])
				}
			}

			if len(ans) == 0 {
				result = make([]string, anLen)
			}

			for k, v := range answers {
				result[k] = v
			}
			return strings.Join(result, "-$-")
		}
	}
	return ""
}

type Audio struct {
	AudioId   string `json:"audioId"`
	Duration  int    `json:"duration"`
	AudioType int    `json:"audioType"`
}

type Option struct {
	Options    []string `json:"options"`
	OptionType int      `json:"optionType"`
}

type Accessories struct {
	Audio  Audio  `json:"audio"`
	Option Option `json:"option"`
}

func makeNewOptions(questionId int64, newOptions map[int]string) string {
	var result string

	//先把答案查出来
	var answers []string
	GetDb().Table("t_questions").Where("F_question_id = ?", questionId).Pluck("F_accessories", &answers)

	if len(answers) > 0 {
		answer := answers[0]
		var accessories Accessories
		jDecoder := json.NewDecoder(strings.NewReader(answer))
		jDecoder.Decode(&accessories)

		if accessories.Option.OptionType != 0 {
			ansJson := make(map[string]interface{})

			for k, v := range newOptions {
				accessories.Option.Options[k] = v
			}

			ansJson["option"] = accessories.Option

			if accessories.Audio.AudioType != 0 {
				ansJson["audio"] = accessories.Audio
			}

			var buffer bytes.Buffer
			jEncoder := json.NewEncoder(&buffer)
			jEncoder.SetEscapeHTML(false)
			jEncoder.Encode(ansJson)

			result = buffer.String()
		}
	}

	return result
}
