package models

import (
	"time"
)

type Project struct {
	ID         int       `gorm:"id" json:"id"`
	Name       string    `gorm:"name" json:"name"`
	Phone      int       `gorm:"phone" json:"phone"`
	Icon       string    `gorm:"icon" json:"icon"`
	Background string    `gorm:"background" json:"background"`
	Type       string    `gorm:"type" json:"type"`
	Address    string    `gorm:"address" json:"address"`
	Money      string    `gorm:"money" json:"money"`
	Agency     string    `gorm:"agency" json:"agency"`
	Introduce  string    `gorm:"introduce" json:"introduce"`
	Status     int       `gorm:"status" json:"status"`
	AddTip     string    `gorm:"column:addtip" json:"addtip"`
	IDCard     string    `gorm:"column:idcard" json:"idcard"`
	Time       time.Time `gorm:"column:time" json:"time"`
}
