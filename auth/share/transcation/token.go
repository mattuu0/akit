package transcation

import (
	"encoding/base64"
	"log"
	"time"

	"github.com/google/uuid"
)

//トークンを保存する
func SaveToken(token string) (string,error) {
	//格納するID生成
	uid,err := uuid.NewRandom()

	//エラー処理
	if err != nil {
		return "",err
	}

	//文字列ID
	uid_str := uid.String()

	//URL Safe URLエンコード
	b64id := base64.URLEncoding.EncodeToString([]byte(uid_str))

	//格納 (有効期限 5分)
	err = Save(uid_str, token,time.Duration(time.Minute * 5))

	//エラー処理
	if err != nil {
		log.Println(err)
		return "",err
	}

	return b64id, nil
}

func GetToken(tokenid string) (string,error) {
	//URLデコード
	decodeid,err := base64.URLEncoding.DecodeString(tokenid)

	//エラー処理
	if err != nil {
		log.Println(err)
		return "",err
	}

	//文字列ID
	uid := string(decodeid)

	//トークン取得
	token,err := Get(uid)

	//エラー処理
	if err != nil {
		log.Println(err)
		return "",err
	}

	return token, nil
}