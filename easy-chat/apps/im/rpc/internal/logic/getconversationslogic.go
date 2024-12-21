package logic

import (
	"context"

	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/im/rpc/internal/svc"
	"easy-chat/pkg/xerr"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话
func (l *GetConversationsLogic) GetConversations(in *im.GetConversationsReq) (*im.GetConversationsResp, error) {
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)

	if err != nil {
		if err == immodels.ErrNotFound {
			return &im.GetConversationsResp{}, nil
		}

		return nil, errors.Wrapf(xerr.NewDBErr(), "find conversation by userId:%s", in.UserId)
	}

	var res im.GetConversationsResp
	copier.Copy(&res, &data)

	ids := make([]string, 0, len(data.ConversationList))
	for _, v := range data.ConversationList {
		ids = append(ids, v.ConversationId)
	}

	conversations, err := l.svcCtx.ConversationModel.ListByConversationIds(l.ctx, ids)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list conversation by conversationIds:%v", ids)
	}

	for _, conversation := range conversations {
		if _, ok := res.ConversationList[conversation.ConversationId]; !ok {
			continue
		}

		total := res.ConversationList[conversation.ConversationId].Total

		if total < int32(conversation.Total) {
			res.ConversationList[conversation.ConversationId].Total = int32(conversation.Total)

			res.ConversationList[conversation.ConversationId].ToRead = int32(conversation.Total) - total

			res.ConversationList[conversation.ConversationId].IsShow = true
		}
	}

	return &res, nil
}
