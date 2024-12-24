package msgTransfer

import (
	"easy-chat/apps/im/ws/ws"
	"easy-chat/pkg/constants"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type groupMsgRead struct {
	mu sync.Mutex

	conversationId string
	push           *ws.Push
	pushCn         chan *ws.Push
	count          int

	pushTime time.Time
	done     chan struct{}
}

func newGroupMsgRead(push *ws.Push, pushCn chan *ws.Push) *groupMsgRead {
	m := &groupMsgRead{
		conversationId: push.ConversationId,
		push:           push,
		pushCn:         pushCn,
		count:          0,
		pushTime:       time.Now(),
		done:           make(chan struct{}),
	}

	go m.transfer()
	return m
}

func (m *groupMsgRead) mergePush(push *ws.Push) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.push == nil {
		m.push = push
	}

	m.count++
	for msgId, read := range push.ReadRecords {
		m.push.ReadRecords[msgId] = read
	}
}

func (m *groupMsgRead) transfer() {
	timer := time.NewTimer(GroupMsgReadRecordDelayTime / 2)
	defer timer.Stop()

	for {
		select {
		case <-m.done:
			return
		case <-timer.C:
			m.mu.Lock()
			pushTime := m.pushTime
			val := GroupMsgReadRecordDelayTime*2 - time.Since(pushTime)
			push := m.push

			if val > 0 && m.count < GroupMsgReadRecordDelayCount || push == nil {
				if val > 0 {
					timer.Reset(val)
				}
				m.mu.Unlock()
				continue
			}

			m.pushTime = time.Now()
			m.push = nil
			m.count = 0
			timer.Reset(GroupMsgReadRecordDelayTime / 2)
			m.mu.Unlock()

			logx.Infof("groupMsgRead transfer push: %v", push)
			m.pushCn <- push
		default:
			m.mu.Lock()

			if m.count >= GroupMsgReadRecordDelayCount {
				push := m.push

				m.push = nil
				m.count = 0
				m.mu.Unlock()

				logx.Infof("groupMsgRead transfer push: %v", push)
				m.pushCn <- push
				continue
			}

			if m.isIdle() {
				m.mu.Unlock()
				m.pushCn <- &ws.Push{
					ChatType:       constants.GroupChatType,
					ConversationId: m.conversationId,
				}
				continue
			}

			m.mu.Unlock()
			tempDelay := GroupMsgReadRecordDelayTime / 4
			if tempDelay > time.Second {
				tempDelay = time.Second
			}

			time.Sleep(tempDelay)
		}
	}

}

func (m *groupMsgRead) IsIdle() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isIdle()
}

func (m *groupMsgRead) isIdle() bool {
	pushTime := m.pushTime
	val := GroupMsgReadRecordDelayTime*2 - time.Since(pushTime)
	if val <= 0 && m.push == nil && m.count == 0 {
		return true
	}

	return false
}

func (m *groupMsgRead) clear() {
	select {
	case <-m.done:
	default:
		close(m.done)
	}

	m.push = nil
}
