package main

import (
	pprof "github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"test.com/project-api/api/midd"
	_ "test.com/project-api/api/project"
	_ "test.com/project-api/api/user"
	"test.com/project-api/config"
	"test.com/project-api/router"
	srv "test.com/project-common"
)

func main() {
	r := gin.Default()
	r.Use(midd.RequestLog())
	// StaticFS 用于将一个文件夹中的文件作为静态文件服务。
	r.StaticFS("/upload", http.Dir("upload"))
	// 注册所有路由
	router.InitRouter(r)
	// 开启pprof 默认的访问路径是 /debug/pprof, 可修改
	pprof.Register(r)
	// 测试代码，模拟内存泄漏；一般内存泄漏大多是goroutine泄漏
	//r.GET("/mem", func(c *gin.Context) {
	//	// 业务代码运行
	//	outCh := make(chan int)
	//	// 每秒起10个goroutine，goroutine会阻塞，不释放内存
	//	tick := time.Tick(time.Second / 10)
	//	i := 0
	//	for range tick {
	//		i++
	//		fmt.Println(i)
	//		alloc1(outCh) // 不停的有goruntine因为outCh堵塞（没人给outch值，无缓冲的chan），无法释放
	//	}
	//})

	srv.Run(r, config.Conf.SC.Name, config.Conf.SC.Addr, nil)
}

//// 一个外层函数
//func alloc1(outCh chan<- int) {
//	go alloc2(outCh)
//}
//
//// 一个内层函数
//func alloc2(outCh chan<- int) {
//	func() {
//		defer fmt.Println("alloc-fm exit")
//		// 分配内存，假用一下
//		buf := make([]byte, 1024*1024*10)
//		_ = len(buf)
//		fmt.Println("alloc done")
//
//		outCh <- 0
//	}()
//}
