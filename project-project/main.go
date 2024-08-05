package main

import (
	"github.com/gin-gonic/gin"
	srv "test.com/project-common"
	"test.com/project-project/config"
	"test.com/project-project/router"

	_ "test.com/project-user/api/user" // 这可以不用要；
)

func main() {
	r := gin.Default()
	// 注册所有路由
	router.InitRouter(r)
	// 注册GRPC
	gc := router.RegisterGrpc()
	// 把GRPC服务注册到ETCD
	router.RegisterEtcdServer()

	stop := func() { // grpc也需要优雅启停；
		gc.Stop()
	}
	srv.Run(r, config.Conf.SC.Name, config.Conf.SC.Addr, stop)
}
