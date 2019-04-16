package models

type IndustryTag struct {
	ID   int    `gorm:"column:id" json:"tagId"`
	Name string `gorm:"column:name" json:"tagName"`
}

type UserIndustry struct {
	Phone    int `gorm:"column:phone" json:"phone"`
	Industry int `gorm:"column:industry" json:"industry"`
}

type ProjectIndustry struct {
	Project  int `gorm:"column:project" json:"project"`
	Industry int `gorm:"column:industry" json:"industry"`
}
