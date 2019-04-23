package models

import (
	"time"
)

type UserInfoSimple struct {
	Phone     int       `gorm:"column:phone" json:"phone"`
	Username  string    `gorm:"column:username" json:"username"`
	Password  string    `gorm:"column:password" json:"password"`
	Icon      string    `gorm:"column:icon" json:"icon"`
	Job       string    `gorm:"column:job" json:"job"`
	Position  string    `gorm:"column:position" json:"position"`
	Profess   string    `gorm:"column:profess" json:"profess"`
	Agency    string    `gorm:"column:agency" json:"agency"`
	Address   string    `gorm:"column:address" json:"address"`
	Introduce string    `gorm:"column:introduce" json:"introduce"`
	School    string    `gorm:"column:school" json:"school"`
	Achieve   string    `gorm:"column:achieve" json:"achieve"`
	Register  time.Time `gorm:"column:register" json:"register"`
}

type UserInfoDetail struct {
	UserInfoSimple UserInfoSimple `json:"userInfoSimple"`
	Connection     Connection     `json:"connection"`
	Tags           []string       `json:"tags"`
}
