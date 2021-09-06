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

// 駅マスタ
type Station struct {
	StationNo   int `gorm:"primary_key"`
	StationName string
	Ido         float64
	Keido       float64
}
