package logic

import (
	"context"

	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/rpc/im"
	"easy-chat/apps/im/rpc/internal/svc"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/wuid"
	"easy-chat/pkg/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SetUpUserConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 建立会话: 群聊, 私聊
func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	switch constants.ChatType(in.ChatType) {
	case constants.SingleChatType:
		conversationId := wuid.CombineId(in.SendId, in.RecvId)

		conversationRes, err := l.svcCtx.ConversationModel.FindOne(l.ctx, conversationId)
		if err != nil {
			if err == immodels.ErrNotFound {
				err := l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
					ConversationId: conversationId,
					ChatType:       constants.SingleChatType,
				})

				if err != nil {
					return nil, errors.Wrapf(xerr.NewDBErr(), "insert conversation by conversationId:%s", conversationId)
				}
			} else {
				return nil, errors.Wrapf(xerr.NewDBErr(), "find conversation by conversationId:%s", conversationId)
			}

		} else if conversationRes != nil {
			return nil, nil
		}

		err = l.setUpUserConversation(conversationId, in.SendId, in.RecvId, constants.SingleChatType, true)

		if err != nil {
			return nil, err
		}

		err = l.setUpUserConversation(conversationId, in.SendId, in.RecvId, constants.SingleChatType, false)
		if err != nil {
			return nil, err
		}

	}

	return &im.SetUpUserConversationResp{}, nil
}

func (l *SetUpUserConversationLogic) setUpUserConversation(conversationId, userId, recvId string,
	chatType constants.ChatType, isShow bool) error {
	// 用户的会话列表
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	if err != nil {
		if err == immodels.ErrNotFound {
			conversations = &immodels.Conversations{
				ID:               primitive.NewObjectID(),
				UserId:           userId,
				ConversationList: make(map[string]*immodels.Conversation),
			}
		} else {
			return errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.FindOne err %v, req %v", err, userId)
		}
	}

	// 更新会话记录
	if _, ok := conversations.ConversationList[conversationId]; ok {
		return nil
	}

	// 添加会话记录
	conversations.ConversationList[conversationId] = &immodels.Conversation{
		ConversationId: conversationId,
		ChatType:       constants.SingleChatType,
		IsShow:         isShow,
	}

	// 更新
	err = l.svcCtx.ConversationsModel.Update(l.ctx, conversations)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "ConversationsModel.Insert err %v, req %v", err, conversations)
	}
	return nil
}