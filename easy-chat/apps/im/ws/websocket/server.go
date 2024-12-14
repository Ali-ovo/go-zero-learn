package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type Server struct {
	sync.RWMutex

	authentication Authentication

	routes map[string]HandlerFunc
	addr   string
	patten string

	connToUser map[*websocket.Conn]string
	userToConn map[string]*websocket.Conn

	upgrader websocket.Upgrader
	logx.Logger
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)

	return &Server{
		routes:   make(map[string]HandlerFunc),
		addr:     addr,
		upgrader: websocket.Upgrader{},

		authentication: opt.Authentication,
		patten:         opt.patten,
		connToUser:     make(map[*websocket.Conn]string),
		userToConn:     make(map[string]*websocket.Conn),

		Logger: logx.WithContext(context.Background()),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover:%v", r)
		}
	}()

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade err: %v", err)
		return
	}

	if !s.authentication.Auth(w, r) {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintln("auth failed")))
		conn.Close()
		return
	}

	s.addCount(conn, r)

	go s.handlerConn(conn)
}

// handlerConn handles websocket connection
func (s *Server) handlerConn(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)

			s.Close(conn)
			return
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshal err %v", err)

			s.Close(conn)
			return
		}

		if handler, ok := s.routes[message.Method]; ok {
			handler(s, conn, &message)
		} else {
			conn.WriteMessage(websocket.TextMessage, []byte("method not found"))
		}
	}
}

func (s *Server) addCount(conn *websocket.Conn, req *http.Request) {
	uid := s.authentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

func (s *Server) GetConn(uid string) *websocket.Conn {
	s.RWMutex.RLock()

	defer s.RWMutex.RUnlock()

	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*websocket.Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*websocket.Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

func (s *Server) GetUsers(conns ...*websocket.Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) Close(conn *websocket.Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		// 已经被关闭
		return
	}

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)

	conn.Close()
}

func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(sendIds...)...)
}

func (s *Server) Send(msg interface{}, conns ...*websocket.Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("Server is stopping")
}
