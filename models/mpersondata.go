package models

import "mime/multipart"

type PersonTag struct {
	ID    int    `gorm:"column:id" json:"tagId"`
	Name  string `gorm:"column:name" json:"tagName"`
	Icon  string `gorm:"column:icon" json:"tagIcon"`
	Plain int    `gorm:"column:plain" json:"plain"`
}

type UserTag struct {
	Phone  int `gorm:"column:phone" json:"phone"`
	Person int `gorm:"column:person" json:"person"`
}

type TagParam struct {
	Name string
	File multipart.File
}
