package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dbconn *gorm.DB = nil
)

func Init() {
	//DB接続
	db,err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})

	//エラー処理
	if err != nil {
		log.Fatalln(err)
	}

	//グローバル変数に格納
	dbconn = db

	//テーブル作成
	err = dbconn.AutoMigrate(
		&Token{}, &User{},
	)

	//エラー処理
	if err != nil {
		log.Fatalln(err)
	}

	//トークン初期化
	TokenInit()
}