package main

import (
	"github.com/gin-gonic/gin"

	"app/auth_grpc/client"
)

func main() {
	

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello World!")
	})

	router.Run(":3001")
}