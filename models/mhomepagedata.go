package models

type HomeShow struct {
	ID      int    `gorm:"column:id" json:"id"`
	Name    string `gorm:"column:name" json:"name"`
	Content string `gorm:"column:content" json:"content"`
	Icon    string `gorm:"column:icon" json:"icon"`
	Type    int    `gorm:"column:type" json:"type"`
	URL     string `gorm:"column:url" json:"url"`
}

type HomeContent struct {
	ID      int    `gorm:"column:id" json:"id"`
	Type    int    `gorm:"column:type" json:"type"`
	Index   int    `gorm:"column:index" json:"index"`
	Content string `gorm:"column:content" json:"content"`
}

type HomeShowDetail struct {
	HomeShow     HomeShow      `json:"homeShow"`
	HomeContents []HomeContent `json:"homeContents"`
}

type HomePageType struct {
	ID   int
	Name string
}

type HomePageContent struct {
	TypeID  int
	Content string
}
