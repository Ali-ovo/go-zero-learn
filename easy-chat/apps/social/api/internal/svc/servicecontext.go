package svc

import (
	"easy-chat/apps/im/rpc/imclient"
	"easy-chat/apps/social/api/internal/config"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/user/rpc/userclient"
	"easy-chat/pkg/interceptor"
	"easy-chat/pkg/middleware"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                config.Config
	IdempotenceMiddleware rest.Middleware
	LimitMiddleware       rest.Middleware

	*redis.Redis
	socialclient.Social
	userclient.User
	imclient.Im
}

// retryPolicy 定义了 gRPC 客户端的重试策略
var retryPolicy = `{
	"methodConfig": [{
		"name": [{
			"service": "social.social"
		}],
		"waitForReady": true,
		"retryPolicy": {
			"maxAttempts": 5,
			"initialBackoff": "0.001s",
			"maxBackoff": "0.002s",
			"backoffMultiplier": 1.0,
			"retryableStatusCodes": ["UNKNOWN", "DEADLINE_EXCEEDED"]
		}
	}]
}`

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		Redis: redis.MustNewRedis(c.Redisx),
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc,
			// zrpc.WithDialOption(grpc.WithDefaultServiceConfig(retryPolicy)),
			zrpc.WithUnaryClientInterceptor(interceptor.DefaultIdempotentClient)),
		),
		User: userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Im:   imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),

		IdempotenceMiddleware: middleware.NewIdempotenceMiddleware().Handler,
		LimitMiddleware:       middleware.NewLimitMiddleware(c.Redisx).TokenLimitHandler(100, 100),
	}
}
