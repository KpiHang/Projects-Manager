package interceptor

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"test.com/project-common/encrypts"
	"test.com/project-grpc/user/login"
	"test.com/project-user/internal/dao"
	"test.com/project-user/internal/repo"
	"time"
)

// CacheInterceptor 除了缓存拦截器，还可以日志拦截器，可以有多个拦截器，就像中间件一样多个；
type CacheInterceptor struct {
	cache    repo.Cache
	cacheMap map[string]any
}

func New() *CacheInterceptor {
	cacheMap := make(map[string]any)
	cacheMap["/login.service.v1.LoginService/MyOrgList"] = &login.OrgListResponse{}
	cacheMap["/login.service.v1.LoginService/FindMemInfoById"] = &login.MemberMessage{}
	return &CacheInterceptor{cache: dao.Rc, cacheMap: cacheMap}
}

type CacheRespOption struct {
	path   string
	typ    any
	expire time.Duration
}

func (c *CacheInterceptor) Cache() grpc.ServerOption {
	return grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 只对当前模块下的指定路径进行拦截；
		respType := c.cacheMap[info.FullMethod]
		if respType == nil { // 说明当前路径不需要拦截；
			return handler(ctx, req)
		}

		// 先查询是否缓存，有的话直接返回；无先请求然后存入缓存；
		con, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		marshal, _ := json.Marshal(req)
		cacheKey := encrypts.Md5(string(marshal))

		respJson, _ := c.cache.Get(con, info.FullMethod+"::"+cacheKey)
		if respJson != "" { // 缓存中有数据，直接返回
			json.Unmarshal([]byte(respJson), &respType)
			zap.L().Info(info.FullMethod + "走了缓存")
			return respType, nil // 这里返回的类型就是grpc往后继续走的 message 类型
		}
		// 缓存出错或者没有缓存，继续请求数据库；
		resp, err = handler(ctx, req) // grpc服务处理请求。请求结果进行缓存
		bytes, _ := json.Marshal(resp)
		c.cache.Put(con, info.FullMethod+"::"+cacheKey, string(bytes), 5*time.Minute)
		zap.L().Info(info.FullMethod + "放入缓存")
		return
	})
}
