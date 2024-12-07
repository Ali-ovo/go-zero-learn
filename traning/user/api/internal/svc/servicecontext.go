package svc

import (
	"go-zero-learn/api/internal/config"
	"go-zero-learn/api/internal/middleware"
	"go-zero-learn/rpc/userclient"

	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config config.Config

	userclient.User

	LoginVerification rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		// User:              userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		LoginVerification: middleware.NewLoginVerificationMiddleware().Handle,
	}
}
