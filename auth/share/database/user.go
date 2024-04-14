package database

import (
	"crypto/sha512"
	"fmt"
)

type User struct {
	//ユーザID
	UserID string `gorm:"primaryKey"`

	//プロバイダのユーザID
	ProviderID string

	//アカウント名
	Name string

	//メールアドレス
	Email string

	//プロバイダ
	Provider string

	//アイコンURL
	IconURL string

	//アイコンパス
	IconPath string
}

// ユーザ作成
func CreateUser(usr User) error {
	//ユーザ作成
	result := dbconn.Create(&usr)

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ユーザ更新
func UpdateUser(usr User) error {
	//ユーザ作成
	result := dbconn.Save(&usr)

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// アイコン取得
func (usr User) GetIcon() string {
	if usr.IconPath != "" {
		return usr.IconPath
	} else {
		return usr.IconURL
	}
}

// ユーザ取得
func GetUser(userID string) (User, error) {
	var user User
	//ユーザ取得
	result := dbconn.First(&user, User{UserID: userID})

	//エラー処理
	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}

// ユーザ削除
func DeleteUSer(userID string) error {
	//ユーザ削除
	result := dbconn.Delete(&User{UserID: userID})

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ID取得
func GetID(provider string, userid string) string {
	//ユーザID
	sha512 := sha512.Sum512([]byte(provider + userid))

	//文字列化
	return fmt.Sprintf("%x", sha512)
}
