package user

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	LoginServiceV1 "test.com/project-user/pkg/service/login.service.v1"
)

var LoginServiceClient LoginServiceV1.LoginServiceClient

func InitRpcUserClient() {
	// grpc.Dial 已被舍弃，使用 grpc.NewClient
	conn, err := grpc.NewClient("127.0.0.1:8881", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	LoginServiceClient = LoginServiceV1.NewLoginServiceClient(conn)
}
