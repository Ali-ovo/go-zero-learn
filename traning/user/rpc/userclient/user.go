// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3
// Source: user.proto

package userclient

import (
	"context"

	"go-zero-learn/traning/user/rpc/user"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CreateReq   = user.CreateReq
	CreateResp  = user.CreateResp
	GetUserReq  = user.GetUserReq
	GetUserResp = user.GetUserResp
	Request     = user.Request
	Response    = user.Response

	User interface {
		// 定义rpc方法
		GetUser(ctx context.Context, in *GetUserReq, opts ...grpc.CallOption) (*GetUserResp, error)
		Create(ctx context.Context, in *CreateReq, opts ...grpc.CallOption) (*CreateResp, error)
		Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	}

	defaultUser struct {
		cli zrpc.Client
	}
)

func NewUser(cli zrpc.Client) User {
	return &defaultUser{
		cli: cli,
	}
}

// 定义rpc方法
func (m *defaultUser) GetUser(ctx context.Context, in *GetUserReq, opts ...grpc.CallOption) (*GetUserResp, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.GetUser(ctx, in, opts...)
}

func (m *defaultUser) Create(ctx context.Context, in *CreateReq, opts ...grpc.CallOption) (*CreateResp, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.Create(ctx, in, opts...)
}

func (m *defaultUser) Ping(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.Ping(ctx, in, opts...)
}
