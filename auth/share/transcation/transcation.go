package transcation

import (
	"context"
	"errors"
	"time"
)

//有効期限付きで保存
func Save(key string, value string,duration time.Duration) error {
	//初期化されているか
	if !isinit {
		return errors.New("not init")
	}

	//Redisに保存
	err := redis_conn.Set(context.Background(), key, value, duration).Err()

	return err
}

//取得
func Get(key string) (string, error) {
	//初期化されているか
	if !isinit {
		return "", errors.New("not init")
	}
	
	//Redisから取得
	result, err := redis_conn.Get(context.Background(), key).Result()

	return result, err
}