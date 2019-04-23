package models

import "time"

type Connection struct {
	Phone  int       `gorm:"column:phone" json:"phone"`
	Card   string    `gorm:"column:card" json:"card"`
	Status int       `gorm:"column:status" json:"status"`
	Look   int       `gorm:"column:look" json:"look"`
	Good   int       `gorm:"column:good" json:"good"`
	Time   time.Time `gorm:"column:time" json:"time"`
	From   int       `gorm:"column:from" json:"from"`
}
