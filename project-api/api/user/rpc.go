package user

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"test.com/project-api/config"
	"test.com/project-common/discovery"
	"test.com/project-common/logs"
	LoginServiceV1 "test.com/project-user/pkg/service/login.service.v1"
)

var LoginServiceClient LoginServiceV1.LoginServiceClient

func InitRpcUserClient() {
	// etcd 相关
	etcdRegister := discovery.NewResolver(config.Conf.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)

	// 从etcd中获取服务；服务发现
	conn, err := grpc.Dial("etcd:///user", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// conn, err := grpc.Dial("127.0.0.1:8881", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	LoginServiceClient = LoginServiceV1.NewLoginServiceClient(conn)
}
