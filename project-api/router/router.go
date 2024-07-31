package router

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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
