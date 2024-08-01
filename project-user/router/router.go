package router

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
	"test.com/project-common/discovery"
	"test.com/project-common/logs"
	"test.com/project-user/config"
	LoginServiceV1 "test.com/project-user/pkg/service/login.service.v1"
)

type Router interface {
	Register(r *gin.Engine)
}

type RegisterRouter struct {
}

func NewRegisterRouter() *RegisterRouter {
	return &RegisterRouter{}
}

func (RegisterRouter) Route(router Router, r *gin.Engine) {
	router.Register(r)
}

var routers []Router

func InitRouter(r *gin.Engine) {
	//router := NewRegisterRouter()
	////以后的模块路由在这进行注册
	//router.Route(&user.RouterUser{}, r)
	for _, router := range routers {
		router.Register(r)
	}
}

// AddRouter 因为 routers 是小写，私有变量
// routers包含了所有业务的路由，某个业务没必要知道其他业务的路由是什么。
func AddRouter(ro ...Router) {
	routers = append(routers, ro...)
}

// gRPCConfig grpc 服务相关
type gRPCConfig struct {
	Addr         string
	RegisterFunc func(*grpc.Server)
}

func RegisterGrpc() *grpc.Server {
	c := gRPCConfig{
		Addr: config.Conf.GC.Addr,
		RegisterFunc: func(g *grpc.Server) {
			LoginServiceV1.RegisterLoginServiceServer(g, LoginServiceV1.NewLoginService()) // 生成代码中提供的函数；
		}}
	s := grpc.NewServer()                 // 创建了一个新的gRPC服务器实例 s。
	c.RegisterFunc(s)                     // 将服务注册到gRPC服务器 s 上
	lis, err := net.Listen("tcp", c.Addr) // 在指定的地址 c.Addr 上创建了一个 TCP 监听器 lis
	if err != nil {
		log.Fatalln("cannot listen:", err)
	}

	go func() { // 放到协程里，看main，如果不放到协程里，main无法向下执行了；
		err = s.Serve(lis) // s.Serve(lis) 会阻塞当前 goroutine，开始接受并处理客户端请求。
		if err != nil {
			log.Fatalln("failed to serve:", err)
		}
	}()
	return s
}

func RegisterEtcdServer() {
	etcdRegister := discovery.NewResolver(config.Conf.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)

	info := discovery.Server{
		Name:    config.Conf.GC.Name,
		Addr:    config.Conf.GC.Addr,
		Version: config.Conf.GC.Version,
		Weight:  config.Conf.GC.Weight,
	}
	r := discovery.NewRegister(config.Conf.EtcdConfig.Addrs, logs.LG)
	_, err := r.Register(info, 2)
	if err != nil {
		log.Fatalln(err)
	}
}
