package auth_grpc

import (
	"log"
	"net"

	"authkit/auth_grpc/agrpc"
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
	log.Print(token)
	return &agrpc.TokenResult{}, nil
}
