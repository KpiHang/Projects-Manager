package router

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
	"test.com/project-common/discovery"
	"test.com/project-common/logs"
	"test.com/project-grpc/account"
	"test.com/project-grpc/auth"
	"test.com/project-grpc/department"
	"test.com/project-grpc/menu"
	"test.com/project-grpc/project"
	"test.com/project-grpc/task"
	"test.com/project-project/config"
	"test.com/project-project/internal/interceptor"
	"test.com/project-project/internal/rpc"
	account_service_v1 "test.com/project-project/pkg/service/account.service.v1"
	auth_service_v1 "test.com/project-project/pkg/service/auth.service.v1"
	department_service_v1 "test.com/project-project/pkg/service/department.service.v1"
	menu_service_v1 "test.com/project-project/pkg/service/menu.service.v1"
	project_service_v1 "test.com/project-project/pkg/service/project.service.v1"
	task_service_v1 "test.com/project-project/pkg/service/task.service.v1"
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
			project.RegisterProjectServiceServer(g, project_service_v1.NewProjectService()) // 注册服务
			task.RegisterTaskServiceServer(g, task_service_v1.NewTaskService())
			account.RegisterAccountServiceServer(g, account_service_v1.New()) // 注册domain
			department.RegisterDepartmentServiceServer(g, department_service_v1.New())
			auth.RegisterAuthServiceServer(g, auth_service_v1.New())
			menu.RegisterMenuServiceServer(g, menu_service_v1.New())
		}}

	s := grpc.NewServer(interceptor.New().Cache()) // 创建了一个新的gRPC服务器实例 s。 // 用了拦截器，可以用多个拦截器
	c.RegisterFunc(s)                              // 将服务注册到gRPC服务器 s 上
	lis, err := net.Listen("tcp", c.Addr)          // 在指定的地址 c.Addr 上创建了一个 TCP 监听器 lis
	if err != nil {
		log.Fatalln("cannot listen:", err)
	}

	go func() { // 放到协程里，看main，如果不放到协程里，main无法向下执行了；
		log.Printf("grpc server started as: %s \n", c.Addr)
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

func InitUserRpc() {
	rpc.InitRpcUserClient()
}
