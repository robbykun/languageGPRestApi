package db

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var (
	db  *gorm.DB
	err error
)

// DB初期化
func Init() {
	fmt.Println("db.Init開始")

	// TODO:環境変数取得
	godotenv.Load(".env.development")
	godotenv.Load()

	// DB接続
	db, err = gorm.Open(os.Getenv("DBMS"), os.Getenv("CONNECT"))

	// リトライ
	count := 0
	if err != nil {
		for {
			if err == nil {
				fmt.Println("")
				break
			}
			fmt.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 180 {
				fmt.Println("")
				panic(err)
			}
			db, err = gorm.Open(os.Getenv("DBMS"), os.Getenv("CONNECT"))
		}
	}

	autoMigration()
}

// DB取得
func GetDB() *gorm.DB {
	return db
}

// DB接続終了
func Close() {
	if err := db.Close(); err != nil {
		panic(err)
	}
}

// マイグレーション
func autoMigration() {
	db.AutoMigrate(&Project{})
	db.AutoMigrate(&Language{})
}
