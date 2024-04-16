package main

import (
	"authkit/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


func main() {
	Init()

	//ルーティング
	router := gin.Default()

	//ミドルウェア設定
	router.Use(Middleware())

	router.POST("/getuser", func(ctx *gin.Context) {
		log.Println(ctx.GetBool("success"))

		//認証されているか
		if ctx.GetBool("success") {
			ctx.JSON(http.StatusOK, gin.H{
				"user":    ctx.MustGet("user"),
				"token":   ctx.GetString("token"),
			})
			return
		}

		ctx.JSON(401, gin.H{
			"message": "unauthorized",
		})
	})

	//トークン更新用
	router.POST("/refresh",RefreshToken)

	//更新確定
	router.POST("/refreshs",SubmitToken)

	//認証エンドポイント
	router.GET("/:provider", provider_auth)

	//コールバックエンドポイント
	router.GET("/:provider/callback", provider_callback)

	router.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

//トークンをを設定する
func SetToken(ctx *gin.Context, token string) {
	//LAX Cookie 1ヶ月
	ctx.SetSameSite(http.SameSiteLaxMode)
	//トークン設定
	ctx.SetCookie("token", token, 2592000, "/", "", true, true)
}