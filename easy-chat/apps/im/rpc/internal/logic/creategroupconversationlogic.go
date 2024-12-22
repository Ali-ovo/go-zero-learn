package logic

import (
	"context"

	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/im/rpc/internal/svc"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupConversationLogic {
	return &CreateGroupConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupConversationLogic) CreateGroupConversation(in *im.CreateGroupConversationReq) (*im.CreateGroupConversationResp, error) {

	res := &im.CreateGroupConversationResp{}

	_, err := l.svcCtx.ConversationModel.FindOne(l.ctx, in.GroupId)
	if err == nil {
		return res, nil
	}

	if err != immodels.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find conversation err:%v", err)
	}

	err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
		ConversationId: in.GroupId,
		ChatType:       constants.GroupChatType,
	})

	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert conversation err:%v", err)
	}

	_, err = NewSetUpUserConversationLogic(l.ctx, l.svcCtx).SetUpUserConversation(&im.SetUpUserConversationReq{
		SendId:   in.CreateId,
		RecvId:   in.GroupId,
		ChatType: int32(constants.GroupChatType),
	})

	return res, err
}
