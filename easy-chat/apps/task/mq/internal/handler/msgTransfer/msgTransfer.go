package msgTransfer

import (
	"context"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/task/mq/internal/svc"
	"easy-chat/pkg/constants"

	"github.com/zeromicro/go-zero/core/logx"
)

type baseMsgTransfer struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBaseMsgTransfer(svc *svc.ServiceContext) *baseMsgTransfer {
	return &baseMsgTransfer{
		svcCtx: svc,
		Logger: logx.WithContext(context.Background()),
	}
}

func (m *baseMsgTransfer) Transfer(ctx context.Context, data *ws.Push) error {
	var err error
	switch data.ChatType {
	case constants.GroupChatType:
		err = m.group(ctx, data)

	case constants.SingleChatType:
		err = m.single(ctx, data)

	}

	return err
}

func (m *baseMsgTransfer) single(ctx context.Context, data *ws.Push) error {
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

func (m *baseMsgTransfer) group(ctx context.Context, data *ws.Push) error {

	users, err := m.svcCtx.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.RecvId,
	})

	if err != nil {
		return err
	}

	data.RecvIds = make([]string, 0, len(users.List))
	for _, members := range users.List {
		if members.UserId == data.SendId {
			continue
		}

		data.RecvIds = append(data.RecvIds, members.UserId)
	}

	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}