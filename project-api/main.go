package main

import (
	"github.com/gin-gonic/gin"
	"test.com/project-api/config"
	"test.com/project-api/router"
	srv "test.com/project-common"

	_ "test.com/project-api/api/user"
)

func main() {
	r := gin.Default()
	// 注册所有路由
	router.InitRouter(r)

	srv.Run(r, config.Conf.SC.Name, config.Conf.SC.Addr, nil)
}
