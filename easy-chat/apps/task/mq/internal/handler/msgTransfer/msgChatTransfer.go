package msgTransfer

import (
	"context"
	"easy-chat/apps/im/immodels"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/social/rpc/socialclient"
	"easy-chat/apps/task/mq/internal/svc"
	"easy-chat/apps/task/mq/mq"
	"easy-chat/pkg/constants"
	"encoding/json"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type MsgChatTransfer struct {
	logx.Logger
	svc *svc.ServiceContext
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svc:    svc,
	}
}

func (m *MsgChatTransfer) Consume(key, value string) error {
	fmt.Println("msgChatTransfer", key, value)

	var (
		data mq.MsgChatTransfer
		ctx  = context.Background()
	)

	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	if err := m.addChatLog(ctx, &data); err != nil {
		return err
	}

	switch data.ChatType {
	case constants.SingleChatType:
		return m.single(&data)
	case constants.GroupChatType:
		return m.group(ctx, &data)
	}

	return nil
}

func (m *MsgChatTransfer) single(data *mq.MsgChatTransfer) error {
	// push message
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

func (m *MsgChatTransfer) group(ctx context.Context, data *mq.MsgChatTransfer) error {

	users, err := m.svc.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
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

	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FormId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, data *mq.MsgChatTransfer) error {
	// 记录消息
	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}
	err := m.svc.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}

	return m.svc.ConversationModel.UpdateMsg(ctx, &chatLog)
}
