package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"app/auth_grpc/client"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//初期化
		ctx.Set("success", false)
		ctx.Set("user", nil)
		ctx.Set("token", "")

		//トークン取得
		token,err := GetToken(ctx)

		//エラー処理
		if err != nil {
			log.Println(err)
			ctx.Next()
			return
		}

		//トークンを検証する
		user,err := client.VerifyToken(token)

		//エラー処理
		if err != nil {
			log.Println(err)
			ctx.Next()
			return
		}

		//データを設定する
		ctx.Set("success", true)
		ctx.Set("user", user)
		ctx.Set("token", token)
	}
}

//トークン取得
func GetToken(ctx *gin.Context) (string, error) {
	//クッキーからトークン取得
	token,err := ctx.Cookie("token")

	//エラー処理
	if err != nil {
		return "", err
	}

	//トークンを返す
	return token, nil
}

//トークンをを設定する
func SetToken(ctx *gin.Context, token string) {
	//LAX Cookie 1ヶ月
	ctx.SetSameSite(http.SameSiteLaxMode)
	//トークン設定
	ctx.SetCookie("token", token, 2592000, "/", "", true, true)
}