package db

// 案件
type Project struct {
	ProjectNo string `gorm:"primary_key"`
	Price     int
	Station   string
}

// プログラミング言語
type Language struct {
	ProjectNo    string `gorm:"primary_key"`
	LanguageType string `gorm:"primary_key"`
}
