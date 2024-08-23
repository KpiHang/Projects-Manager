package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test.com/project-api/api/midd"
	"test.com/project-api/config"
	"test.com/project-api/router"
	srv "test.com/project-common"

	_ "test.com/project-api/api/project"
	_ "test.com/project-api/api/user"
)

func main() {
	r := gin.Default()
	r.Use(midd.RequestLog())
	// StaticFS 用于将一个文件夹中的文件作为静态文件服务。
	r.StaticFS("/upload", http.Dir("upload"))
	// 注册所有路由
	router.InitRouter(r)

	srv.Run(r, config.Conf.SC.Name, config.Conf.SC.Addr, nil)
}
