package main

import (
	"authkit/database"
	"context"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"time"
)

// 認証後の処理
func provider_callback(ctx *gin.Context) {
	//プロバイダ取得
	provider := ctx.Param("provider")
	ctx.Request = contextWithProviderName(ctx, provider)

	//認証完了
	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{"message": "failed"})
		return
	}

	//ID取得
	uid := database.GetID(user.Provider, user.UserID)

	//ユーザ取得
	usr, err := database.GetUser(uid)

	//エラー処理
	if err == gorm.ErrRecordNotFound {
		//見つからないとき
		//ユーザ作成
		err = database.CreateUser(database.User{
			UserID:     uid,
			ProviderID: user.UserID,
			Name:       user.Name,
			Email:      user.Email,
			Provider:   user.Provider,
			IconURL:    user.AvatarURL,
			IconPath:   user.AvatarURL,
		})

		//エラー処理
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusOK, gin.H{"message": "failed"})
			return
		}

		//ユーザ取得
		usr, err = database.GetUser(uid)

		//エラー処理
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusOK, gin.H{"message": "failed"})
			return
		}
	} else if err != nil {
		//それ以外のエラー
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{"message": "failed"})
		return
	}

	//ユーザ更新
	err = database.UpdateUser(database.User{
		UserID:     uid,
		ProviderID: user.UserID,
		Name:       user.Name,
		Email:      user.Email,
		Provider:   user.Provider,
		IconURL:    user.AvatarURL,
		IconPath:   user.AvatarURL,
	})

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{"message": "failed"})
		return
	}

	//トークン生成
	token,err := database.GenToken(database.Token{
		UserID: usr.UserID, 
		TokenID: database.GenID(), 
		BaseID: "",
		UserAgent: ctx.Request.UserAgent(),
		Exptime: time.Now().AddDate(1,0,0),
	})

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{"message": "failed"})
		return
	}

	//LAX Cookie
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("token", token, 3600000, "/", "", false, true)
	//リダイレクト
	ctx.Redirect(http.StatusFound, "/auth/ping")
}

// 認証
func provider_auth(ctx *gin.Context) {
	//プロバイダ設定
	provider := ctx.Param("provider")
	ctx.Request = contextWithProviderName(ctx, provider)

	//認証開始
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

// プロバイダ取得
func contextWithProviderName(ctx *gin.Context, provider string) *http.Request {
	return ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", provider))
}
