package main

import (
	"authkit/database"
	"authkit/transcation"
	"context"
	"errors"

	"github.com/gin-contrib/sessions"
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to complete user auth"}) 
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
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create user"})
			return
		}

		//ユーザ取得
		usr, err = database.GetUser(uid)

		//エラー処理
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get user"})
			return
		}
	} else if err != nil {
		//それ以外のエラー
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get user"})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update user"})
		return
	}

	//トークン生成
	token,err := database.GenToken(database.Token{
		UserID: usr.UserID, 
		TokenID: database.GenID(), 
		BaseID: "",
		UserAgent: ctx.Request.UserAgent(),
		Exptime: time.Now().AddDate(0,1,0),	//有効期限1ヶ月
	})

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{"message": "failed to generate token"})
		return
	}

	log.Println(user.Name)
	log.Println(user.FirstName)
	log.Println(user.LastName)
	log.Println(user.NickName)

	//認証トークンを設定
	tokenid,err := transcation.SaveToken(token)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to save token"})
		return
	}

	//リダイレクトURL取得
	redirect_url,err := GetRedirect_URL(ctx)

	log.Println(redirect_url)
	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{"message": "invalid redirect_url"})
		return
	}

	//リダイレクト
	ctx.Redirect(http.StatusFound, redirect_url + "?token=" + tokenid)
}

// 認証
func provider_auth(ctx *gin.Context) {
	//リダイレクトURL設定
	err := SetRedirect_URL(ctx)

	//エラー処理
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusOK, gin.H{"message": "failed to set redirect_url"})
		return
	}

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

//トークン更新
func RefreshToken(ctx *gin.Context) {
	//認証されているか
	if !ctx.GetBool("success") {
		//されていないとき
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	//トークン取得
	token,exits := ctx.Get("token")

	//トークン無し
	if !exits {
		//エラーを返す
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	//トークン解析
	token_data,_ := database.ValidToken(token.(string))

	//新しいトークン取得
	new_token,err := database.UpdateToken(token_data)

	//エラー処理
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed",
		})
		return
	}

	//トークン設定
	SetToken(ctx, new_token)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

//トークン更新確定
func SubmitToken(ctx *gin.Context) {
	//認証されているか
	if !ctx.GetBool("success") {
		//されていないとき
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	//トークン取得
	token,exits := ctx.Get("token")

	//トークン無し
	if !exits {
		//エラーを返す
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	//トークン解析
	token_data,_ := database.ValidToken(token.(string))

	//更新確定
	err := database.SubmitUpdate(token_data)

	//エラー処理
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func SetRedirect_URL(ctx *gin.Context) (error) {
	//リダイレクトURL取得
	redirect_url := ctx.DefaultQuery("redirect_url","")

	//URL検証
	if redirect_url == "" {
		return errors.New("redirect_url is required")
	}

	//セッション取得
	session := sessions.Default(ctx)
	//URL格納
	session.Set("redirect_url", redirect_url)
	err := session.Save()

	//エラー処理
	if err != nil {
		return err
	}

	return nil
}


func GetRedirect_URL(ctx *gin.Context) (string,error) {
	//セッション取得
	session := sessions.Default(ctx)

	//リダイレクトURL取得
	redirect_url := session.Get("redirect_url")

	//存在するか
	if redirect_url == nil {
		return "",errors.New("redirect_url is not found")
	}

	return redirect_url.(string), nil
}