package msgTransfer

import (
	"easy-chat/apps/im/ws/ws"
	"sync"
	"time"
)

type groupMsgRead struct {
	mu sync.Mutex

	push   *ws.Push
	pushCn chan *ws.Push
	count  int

	pushTime time.Time
	done     chan struct{}
}

func newGroupMsgRead(push *ws.Push, pushCn chan *ws.Push) *groupMsgRead {
	return &groupMsgRead{
		push:     push,
		pushCn:   pushCn,
		count:    0,
		pushTime: time.Now(),
		done:     make(chan struct{}),
	}
}
