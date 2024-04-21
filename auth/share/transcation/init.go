package transcation

import (
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

//グローバル変数
var (
	isinit = false
	redis_conn *redis.Client = nil
)

//文字列を数字に変える
func String_To_Int(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Init() {
	//Redis接続
	redis_conn = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("Redis_Host"), os.Getenv("Redis_Port")),
		Password: os.Getenv("Redis_Password"),
		DB:       String_To_Int(os.Getenv("Transaction_Redis_DB")),
	})

	//初期化済みにする
	isinit = true
}