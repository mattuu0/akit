package auth_grpc

import (
	"errors"
	"log"
	"net"
	"os"

	"authkit/auth_grpc/agrpc"
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
func (auths *AuthService) Logout(context.Context, *agrpc.LogoutToken) (*agrpc.LogoutResult, error) {
	panic("unimplemented")
}

// Verify implements agrpc.AuthServiceServer.
func (auths *AuthService) Verify(context.Context, *agrpc.VerifyToken) (*agrpc.VerifyResult, error) {
	panic("unimplemented")
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
