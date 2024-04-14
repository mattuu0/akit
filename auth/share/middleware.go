package main

import (
	"authkit/database"
	"log"

	"github.com/gin-gonic/gin"
)

//ミドルウェア
func Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//初期化
		ctx.Set("success", false)
		ctx.Set("user",database.User{})
		ctx.Set("token","")

		//トークン取得
		token_str,err := GetToken(ctx)

		//エラー処理
		if err != nil {
			log.Println(err)
			ctx.Next()
			return
		}

		//トークン検証
		token_data,valid := database.ValidToken(token_str)

		//認証できたか
		if !valid {
			//認証できていない場合
			ctx.Next()
			return
		}

		//ユーザー取得
		user,err := database.GetUser(token_data.UserID)

		//エラー処理
		if err != nil {
			log.Println(err)
			ctx.Next()
			return
		}

		//ユーザー格納
		ctx.Set("user",user)
		//成功
		ctx.Set("success", true)
		//トークン格納
		ctx.Set("token",token_str)

		ctx.Next()
	}
}

func GetToken(ctx *gin.Context) (string,error) {
	//トークン取得
	token,err := ctx.Cookie("token")

	//エラー処理
	if err != nil {
		return "", err
	}

	return token, nil
}