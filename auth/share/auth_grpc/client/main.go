package client

import (
	"app/auth_grpc/agrpc"
	"errors"
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	client agrpc.AuthServiceClient = nil
	secret string = ""
	grpc_conn *grpc.ClientConn = nil
)

func Init(Secret string) (error) {
	//コネクション確立
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %s", err)
		return err
	}

	//GRPCクライアント取得
	client = agrpc.NewAuthServiceClient(conn)
	//接続を保存
	grpc_conn = conn
	//認証キーを保存
	secret = Secret

	return nil
}

//トークンを取得する関数
func GetToken(token string) (string, error) {
	//トークンを取得
	result, err := client.GetToken(context.Background(), &agrpc.GetData{
		Secret: os.Getenv("Token_Secret"),
		Token: token,
	})

	if err != nil {
		return "", err
	}

	//取得に失敗したとき
    if !result.Success {
		return "", errors.New("failed to get token")
	}

	return result.Token, nil
}

//認証を確認する関数
func VerifyToken(token string) (*agrpc.User, error) {
	//トークンを検証
	result, err := client.Verify(context.Background(), &agrpc.VerifyToken{
		Token: token,
	})

	//エラー処理
	if err != nil {
		return nil, err
	}

	//検証に失敗したとき
	if !result.Success {
		return nil, errors.New("failed to verify token")
	}

	return result.User, nil
}

//ログアウト
func Logout(token string) error {
	//トークンを削除
	result, err := client.Logout(context.Background(), &agrpc.LogoutToken{
		Secret: os.Getenv("Token_Secret"),
		Token: token,
	})

	//エラー処理
	if err != nil {
		return err
	}

	//削除に失敗したとき
	if !result.Success {
		return errors.New("failed to logout")
	}
	
	return nil
}