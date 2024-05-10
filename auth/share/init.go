package main

import (
	"authkit/auth_grpc/server"
	"authkit/database"
	"authkit/transcation"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	//"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/microsoftonline"

	"github.com/gorilla/sessions"
)

//グローバル変数
var (

)

func Init() {
	//envロード
	loadEnv()

	//認証初期化
	gothic_init()

	//プロバイダ初期化
	Provider_init()

	//DB初期化
	database.Init()

	//GRPC初期化
	go server.Init()
}

//プロバイダ初期化
func Provider_init() {
	goth.UseProviders(
		//discord.New(os.Getenv("DISCORD_CLIENT_ID"), os.Getenv("DISCORD_CLIENT_SECRET"), os.Getenv("DISCORD_CALLBACK_URL"), discord.ScopeIdentify, discord.ScopeEmail),
		microsoftonline.New(os.Getenv("Microsoft_ClientID"), os.Getenv("Microsoft_ClientSecret"), os.Getenv("Microsoft_CallbackURL")),
	)
}

//認証ライブラリ初期化
func gothic_init() {
	session_key := os.Getenv("Session_Key")
	maxAge := os.Getenv("Session_MaxAge")
	ispod := false

	//数字に変換
	max_age, err := strconv.Atoi(maxAge)

	//エラー処理
	if err != nil {
		log.Fatalln("Session_MaxAge is not number")
	}

	//セッションストア
	store := sessions.NewCookieStore([]byte(session_key))
	store.MaxAge(max_age)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = ispod

	//認証設定
	gothic.Store = store

	//トランザクション初期化
	transcation.Init()
}

//環境変数読み込み
func loadEnv() {
	//ファイル読み込み
	err := godotenv.Load(".env")

	//エラー処理
	if err != nil {
		log.Fatalln("cannot load env")
	}
}