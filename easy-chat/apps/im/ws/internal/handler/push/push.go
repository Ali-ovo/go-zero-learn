package push

import (
	"easy-chat/apps/im/ws/internal/svc"
	"easy-chat/apps/im/ws/websocket"
	"easy-chat/apps/im/ws/ws"

	"github.com/mitchellh/mapstructure"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Push

		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			srv.Send(websocket.NewErrMessage(err))

			return
		}

		// send
		rconn := srv.GetConn(data.RecvId)

		if rconn == nil {
			// todo

			return
		}

		srv.Infof("push message to %s", data.RecvId)

		srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendTime:       data.SendTime,
			Msg: ws.Msg{
				MType:   data.MType,
				Content: data.Content,
			},
		}), rconn)
	}
}
