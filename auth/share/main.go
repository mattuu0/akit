package main

import (
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

	router.GET("/ping", func(ctx *gin.Context) {
		log.Println(ctx.GetBool("success"))

		if ctx.GetBool("success") {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "pong",
				"user":    ctx.MustGet("user"),
				"token":   ctx.GetString("token"),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//認証エンドポイント
	router.GET("/:provider", provider_auth)

	//コールバックエンドポイント
	router.GET("/:provider/callback", provider_callback)

	router.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
