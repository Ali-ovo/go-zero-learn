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

	connToUser map[*websocket.Conn]string
	userToConn map[string]*websocket.Conn

	upgrader websocket.Upgrader
	logx.Logger
}

func NewServer(addr string) *Server {
	return &Server{
		routes:   make(map[string]HandlerFunc),
		addr:     addr,
		upgrader: websocket.Upgrader{},

		authentication: new(authentication),
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

	go s.handlerConn(conn)
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

// handlerConn handles websocket connection
func (s *Server) handlerConn(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)

			// todo: close conn
			return
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshal err %v", err)

			return
		}

		if handler, ok := s.routes[message.Method]; ok {
			handler(s, conn, &message)
		} else {
			conn.WriteMessage(websocket.TextMessage, []byte("method not found"))
		}
	}
}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

func (s *Server) Start() {
	http.HandleFunc("/ws", s.ServerWs)
	fmt.Println("Server is running on port 8080")
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("Server is stopping")
}
