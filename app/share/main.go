package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"app/auth_grpc/client"

	"app/auth_grpc/agrpc"
)

func main() {
	//env 読み込み
	loadEnv()

	//GRPCクライアント初期化
	client.Init(os.Getenv("Token_Secret"))

	router := gin.Default()

	//ミドルウェア設定
	router.Use(AuthMiddleware())

	router.GET("/callback", func(ctx *gin.Context) {
		//トークンを取得
		token := ctx.DefaultQuery("token", "")

		log.Println("token: " + token)
		//トークンを取得(引き換え)する
		get_token,err := client.GetToken(token)

		//エラー処理
		if err != nil {
			log.Println(err)
			ctx.JSON(500, gin.H{"error": "failed to get token"})
			return
		}

		//トークンをクッキーに設定
		SetToken(ctx, get_token)		

		ctx.Redirect(http.StatusFound, "/statics/index.html")
	})

	//更新用エンドポイント
	router.POST("/refresh", func(ctx *gin.Context) {
		//認証されているか
		if !ctx.GetBool("success") {
			//認証されていない場合
			ctx.JSON(401, gin.H{"error": "not authenticated"})
			return
		}

		//トークンを取得
		token := ctx.GetString("token")

		//トークンを更新する
		get_token,err := client.RefreshToken(token)

		//エラー処理
		if err != nil {
			log.Println(err)
			ctx.JSON(500, gin.H{"error": "failed to refresh token"})
			return
		}

		//新しいトークンを設定する
		SetToken(ctx, get_token)

		ctx.JSON(200, gin.H{"success": true})
	})

	//更新を確定するエンドポイント
	router.POST("/refreshs", func(ctx *gin.Context) {
		//認証されているか
		if !ctx.GetBool("success") {
			//認証されていない場合
			ctx.JSON(401, gin.H{"error": "not authenticated"})
			return
		}

		//トークンを取得
		token := ctx.GetString("token")

		//トークンを更新する
		err := client.RefreshTokenS(token)

		//エラー処理
		if err != nil {
			log.Println(err)
			ctx.JSON(500, gin.H{"error": "failed to refresh token"})
			return
		}

		ctx.JSON(200, gin.H{"success": true})
	})

	//ログアウトエンドポイント
	router.POST("/logout", func(ctx *gin.Context) {
		//認証されているか
		if !ctx.GetBool("success") {
			//認証されていない場合
			ctx.JSON(401, gin.H{"error": "not authenticated"})
			return
		}

		//ログアウト
		err := client.Logout(ctx.GetString("token"))

		//エラー処理
		if err != nil {
			log.Println(err)
			ctx.JSON(500, gin.H{"error": "failed to logout"})
			return
		}

		//トークンをクッキーから削除
		SetToken(ctx, "")

		ctx.JSON(200, gin.H{"success": true})
	})

	router.GET("/getuser", func(ctx *gin.Context) {
		//認証されているか
		if !ctx.GetBool("success") {
			//認証されていない場合
			ctx.JSON(401, gin.H{"error": "not authenticated"})
			return
		}

		//ユーザー情報を取得
		user,exits := ctx.Get("user")

		//ユーザー情報を取得出来なかった場合
		if !exits {
			ctx.JSON(500, gin.H{"error": "failed to get user"})
			return
		}

		//User型にキャストする
		castuser := user.(*agrpc.User)

		ctx.JSON(200, gin.H{"user": User{
			UserId: castuser.UserId,
			Name: castuser.Name,
			Email: castuser.Email,
			Icon: castuser.Icon,
			Provider: castuser.Provider,
			ProviderID: castuser.ProviderUID,
		}})
		
	})
	router.Run(":3001")
}

type User struct {
	UserId string `json:"user_id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Icon string `json:"icon"`
	Provider string `json:"provider"`
	ProviderID string `json:"provider_id"`
}

// .envを呼び出します。
func loadEnv() {
	// ここで.envファイル全体を読み込みます。
	// この読み込み処理がないと、個々の環境変数が取得出来ません。
	// 読み込めなかったら err にエラーが入ります。
	err := godotenv.Load(".env")
	
	// もし err がnilではないなら、"読み込み出来ませんでした"が出力されます。
	if err != nil {
		log.Fatalf("読み込み出来ませんでした: %v", err)
	} 
}