package client

import (
	"authkit/auth_grpc/agrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	client := agrpc.NewAuthServiceClient(conn)

	response, err := client.GetToken(context.Background(), &agrpc.GetData{
		Token: "token",
		Secret: "secret",
	})

	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Print(response.Success)

	defer conn.Close()
}
