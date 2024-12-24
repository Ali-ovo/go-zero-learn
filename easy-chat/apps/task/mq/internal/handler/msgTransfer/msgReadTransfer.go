package msgTransfer

import (
	"context"
	"easy-chat/apps/im/ws/ws"
	"easy-chat/apps/task/mq/internal/svc"
	"easy-chat/apps/task/mq/mq"
	"easy-chat/pkg/constants"
	"easy-chat/pkg/constants/bitmap"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"
)

var (
	GroupMsgReadRecordDelayTime  = time.Second
	GroupMsgReadRecordDelayCount = 10
)

const (
	GroupMsgReadHandlerAtTransfer = iota
	GroupMsgReadHandlerDelayTransfer
)

type MsgReadTransfer struct {
	*baseMsgTransfer

	cache.Cache
	mu        sync.Mutex
	groupMsgs map[string]*groupMsgRead
	push      chan *ws.Push
}

func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	m := &MsgReadTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svc),
		groupMsgs:       make(map[string]*groupMsgRead, 1),
		push:            make(chan *ws.Push, 1),
	}

	if svc.Config.MsgReadHandler.GroupMsgReadHandler != GroupMsgReadHandlerAtTransfer {
		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount > 0 {
			GroupMsgReadRecordDelayCount = svc.Config.MsgReadHandler.GroupMsgReadRecordDelayCount
		}

		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime > 0 {
			GroupMsgReadRecordDelayTime = time.Duration(svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime) * time.Second
		}
	}

	go m.transfer()

	return m
}

func (m *MsgReadTransfer) Consume(key, value string) error {
	m.Info("MsgReadTransfer", "key", key, "value", value)
	var (
		data mq.MsgMarkRead
		ctx  = context.Background()
	)

	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	readRecords, err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}

	push := &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ContentType:    constants.ContentMarkRead,
		ReadRecords:    readRecords,
		// RecvIds:        data.RecvIds,
		// SendTime:       data.SendTime,
		// MType:          data.MType,
		// Content:        data.Content,
	}

	switch data.ChatType {
	case constants.SingleChatType:
		m.push <- push

	case constants.GroupChatType:
		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			m.push <- push
		}

		m.mu.Lock()
		defer m.mu.Unlock()

		push.SendId = ""
		if _, ok := m.groupMsgs[push.ConversationId]; ok {
			m.Infof("merge push:%v", push)

			m.groupMsgs[push.ConversationId].mergePush(push)
		} else {
			m.Infof("new group push:%v", push)

			m.groupMsgs[push.ConversationId] = newGroupMsgRead(push, m.push)
		}

	}

	// return m.Transfer(ctx,)
	return nil
}

func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {
	res := make(map[string]string)

	chatLogs, err := m.svcCtx.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
	if err != nil {
		return nil, err
	}

	for _, chatLog := range chatLogs {
		switch chatLog.ChatType {
		case constants.SingleChatType:
			chatLog.ReadRecords = []byte{1}
		case constants.GroupChatType:
			readRecords := bitmap.Load(chatLog.ReadRecords)
			readRecords.Set(data.SendId)
			chatLog.ReadRecords = readRecords.Export()
		}

		res[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)

		err = m.svcCtx.ChatLogModel.UpdateMakeRead(ctx, chatLog.ID, chatLog.ReadRecords)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (m *MsgReadTransfer) transfer() {
	for push := range m.push {
		if push.RecvId != "" || len(push.RecvIds) > 0 {
			if err := m.Transfer(context.Background(), push); err != nil {
				m.Errorf("transfer push error:%v", err)
			}
		}

		if push.ChatType == constants.SingleChatType {
			continue
		}

		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			continue
		}

		m.mu.Lock()

		if _, ok := m.groupMsgs[push.ConversationId]; ok && m.groupMsgs[push.ConversationId].IsIdle() {
			m.groupMsgs[push.ConversationId].clear()
			delete(m.groupMsgs, push.ConversationId)
		}

		m.mu.Unlock()
	}
}
