package auth_grpc

import (
	"errors"
	"log"
	"net"
	"os"

	"authkit/auth_grpc/agrpc"
	"authkit/database"
	"authkit/transcation"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	log.Print("main start")

	// 9000番ポートでクライアントからのリクエストを受け付けるようにする
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Sample構造体のアドレスを渡すことで、クライアントからGetDataリクエストされると
	// GetDataメソッドが呼ばれるようになる
	agrpc.RegisterAuthServiceServer(grpcServer, &AuthService{})

	// 以下でリッスンし続ける
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	log.Print("main end")
}

type AuthService struct {
	name string
}

// Logout implements agrpc.AuthServiceServer.
func (auths *AuthService) Logout(ctx context.Context, token_data *agrpc.LogoutToken) (*agrpc.LogoutResult, error) {
	//シークレット検証
	if token_data.Secret != os.Getenv("Token_Secret") {
		//シークレットが一致しない場合
		log.Println("invalid secret")
		return &agrpc.LogoutResult{
			Success: false,
		}, errors.New("invalid secret")
	}

	//トークンを検証する
	valid_data,isvalid := database.ValidToken(token_data.Token)

	//エラー処理
	if !isvalid {
		log.Println("invalid token")
		return &agrpc.LogoutResult{
			Success: false,
		}, errors.New("invalid token")
	}

	//トークンを削除する
	err := database.DeleteToken(valid_data.TokenID)

	//エラー処理
	if err != nil {
		log.Println(err)
		return &agrpc.LogoutResult{
			Success: false,
		}, err
	}

	return &agrpc.LogoutResult{
		Success: true,
	}, nil
}

// Verify implements agrpc.AuthServiceServer.
func (auths *AuthService) Verify(ctx context.Context,token_data *agrpc.VerifyToken) (*agrpc.VerifyResult, error) {
	//トークン検証
	valid_data,isvalid := database.ValidToken(token_data.Token)

	//エラー処理
	if !isvalid {
		log.Println("invalid token")
		return &agrpc.VerifyResult{
			Success: false,
		}, errors.New("invalid token")
	}

	//ユーザ取得
	user_data,err := database.GetUser(valid_data.UserID)

	//エラー処理
	if err != nil {
		log.Println(err)
		return &agrpc.VerifyResult{
			Success: false,
		}, err
	}

	return &agrpc.VerifyResult{
		Success: true,
		User:    &agrpc.User{
			UserId:   valid_data.UserID,
			Name:     user_data.Name,
			Email:    user_data.Email,
			Icon:     user_data.IconURL,
			Provider: user_data.Provider,
			ProviderUID: user_data.ProviderID,
		},
	},nil
}

func (auths *AuthService) GetToken(
	ctx context.Context,
	token *agrpc.GetData,
) (*agrpc.TokenResult, error) {
	//シークレットを検証する
	if token.Secret != os.Getenv("Token_Secret") {
		//トークンが一致しない場合
		log.Println("invalid secret")
		return &agrpc.TokenResult{
			Success: false,
		}, errors.New("invalid secret")
	}

	//トークンを取得する
	get_token,err := transcation.GetToken(token.Token)

	//エラー処理
	if err != nil {
		log.Println(err)
		return &agrpc.TokenResult{
			Success: false,
		}, err
	}

	return &agrpc.TokenResult{
		Success: true,
		Token:   get_token,
	}, nil
}
