package logic

import (
	"context"

	"go-zero-learn/traning/user/api/internal/svc"
	"go-zero-learn/traning/user/api/internal/types"
	"go-zero-learn/traning/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserinfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserinfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserinfoLogic {
	return &UserinfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserinfoLogic) Userinfo(req *types.UserReq) (resp *types.UserResp, err error) {
	getUserResp, err := l.svcCtx.User.GetUser(l.ctx, &userclient.GetUserReq{
		Id: req.Id,
	})

	if err != nil {
		return nil, err
	}

	return &types.UserResp{
		Id:    getUserResp.Id,
		Name:  getUserResp.Name,
		Phone: getUserResp.Phone,
	}, nil
}
