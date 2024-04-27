package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/google/uuid"
)

var (
	redis_conn  *redis.Client = nil
	sign_method               = jwt.SigningMethodHS512
)

// 文字列を数字に変える
func String_To_Int(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func TokenInit() {
	//Redis接続
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("Redis_Host"), os.Getenv("Redis_Port")),
		Password: os.Getenv("Redis_Password"),                // no password set
		DB:       String_To_Int(os.Getenv("Token_Redis_DB")), // use default DB
	})

	//グローバル変数に格納
	redis_conn = rdb

	go deleteExpired()
}

type Token struct {
	TokenID   string `gorm:"primaryKey"`
	UserID    string
	UserAgent string
	BaseID    string
	Exptime   time.Time
}

// エンコード
func (token Token) Encode() (string, error) {
	//マーシャル
	bin, err := msgpack.Marshal(token)

	//エラー処理
	if err != nil {
		return "", err
	}

	return string(bin), nil
}

// デコード
func DecodeToken(token_str string) (Token, error) {
	//デコード
	var token Token
	err := msgpack.Unmarshal([]byte(token_str), &token)

	//エラー処理
	if err != nil {
		return Token{}, err
	}

	return token, nil
}

func GetToken(tokenID string) (Token, error) {
	//コンテキスト
	ctx := context.Background()

	//トークン取得
	result, err := redis_conn.Get(ctx, tokenID).Result()

	//エラー処理
	if err != nil {
		//取得できなかったらDBから取得

		var token Token
		//ユーザ取得
		result := dbconn.First(&token, Token{TokenID: tokenID})

		//エラー処理
		if result.Error != nil {
			return Token{}, result.Error
		}

		//エンコード
		token_str, err := token.Encode()

		//エラー処理
		if err != nil {
			return Token{}, err
		}

		//キャッシュにセット (有効期限 1時間)
		err = redis_conn.Set(ctx, tokenID, token_str, time.Duration(time.Hour*1)).Err()

		//エラー処理
		if err != nil {
			return Token{}, err
		}

		return token, err
	}

	//トークンデコード
	token, err := DecodeToken(result)

	//エラー処理
	if err != nil {
		return Token{}, err
	}

	return token, nil
}

func DeleteToken(tokenID string) error {
	//トークン削除
	result := dbconn.Delete(&Token{TokenID: tokenID})

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// トークン登録
func RegisterToken(token Token) error {
	//トークン作成
	result := dbconn.Save(&token)

	//エラー処理
	if result.Error != nil {
		return result.Error
	}

	//エンコード
	token_str, err := token.Encode()

	//エラー処理
	if err != nil {
		return err
	}

	//キャッシュにセット (有効期限 1時間)
	err = redis_conn.Set(context.Background(), token.TokenID, token_str, time.Duration(time.Hour*1)).Err()

	//エラー処理
	if err != nil {
		return err
	}

	return nil
}

// トークン生成
func GenToken(token_data Token) (string, error) {
	//トークン作成
	claims := jwt.MapClaims{
		"userid":  token_data.UserID,
		"tokenid": token_data.TokenID,
		"baseid":  token_data.BaseID,
	}

	//トークン生成
	token_str := jwt.NewWithClaims(sign_method, claims)

	//トークン署名
	token, err := token_str.SignedString([]byte(os.Getenv("JWT_Secret")))

	//エラー処理
	if err != nil {
		return "", err
	}

	//トークン登録
	err = RegisterToken(token_data)

	//エラー処理
	if err != nil {
		return "", err
	}

	return token, err
}

func ValidToken(tokenString string) (Token, bool) {
	//トークン検証
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_Secret")), nil
	})

	//エラー処理
	if err != nil {
		log.Println(err)
		return Token{}, false
	}

	//トークン検証
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//トークンID
		tokenID := claims["tokenid"].(string)
		//トークン取得
		token, err := GetToken(tokenID)

		//エラー処理
		if err != nil {
			return Token{}, false
		}

		return token, true
	}

	return Token{}, false
}

// ID生成
func GenID() string {
	//ID生成
	uid, err := uuid.NewRandom()

	//エラー処理
	if err != nil {
		log.Println(err)
		return ""
	}

	return uid.String()
}

// 更新用トークン発行
func UpdateToken(base_token Token) (string, error) {
	//有効期限5分のトークン発行
	new_token, err := GenToken(Token{
		UserID:    base_token.UserID,
		TokenID:   GenID(),
		BaseID:    base_token.TokenID,
		Exptime:   time.Now().Add(time.Minute * 5),
		UserAgent: base_token.UserAgent,
	})

	//エラー処理
	if err != nil {
		log.Println(err)
		return "", err
	}

	return new_token, nil
}

func SubmitUpdate(token Token) error {
	//古いトークン削除
	err := DeleteToken(token.BaseID)

	//エラー処理
	if err != nil {
		log.Println(err)
		return err
	}

	//有効期限1ヶ月
	token.Exptime = time.Now().AddDate(0, 1, 0)
	//新しいトークンの有効期限更新
	err = RegisterToken(token)

	//エラー処理
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// 期限切れトークン削除
func deleteExpired() error {
	for {
		//10秒スリープ
		time.Sleep(time.Second * 3)

		//トークン取得
		var tokens []Token
		result := dbconn.Where("Exptime < ?", time.Now()).Find(&tokens)

		//エラー処理
		if result.Error != nil {
			log.Println("delete error", result.Error)
			continue
		}

		//for で回す
		for i := 0; i < len(tokens); i++ {
			//トークン削除
			err := DeleteToken(tokens[i].TokenID)
			//エラー処理
			if err != nil {
				log.Println("delete error", result.Error)
				continue
			}
		}

		if result.RowsAffected != 0 {
			log.Println("delete rows", result.RowsAffected)
		}
	}
}
