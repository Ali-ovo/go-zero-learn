package websocket

import (
	"encoding/json"
	"net/url"

	"github.com/gorilla/websocket"
)

type Client interface {
	Close() error

	Send(v any) error
	Read(v any) error
}

type client struct {
	*websocket.Conn

	host string

	opt dailOption
}

func NewClient(host string, opts ...DailOptions) *client {
	opt := newDailOptions(opts...)

	c := client{
		Conn: nil,
		host: host,
		opt:  opt,
	}

	con, err := c.dail()
	if err != nil {
		panic(err)
	}

	c.Conn = con

	return &c
}

func (c *client) dail() (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: c.host, Path: c.opt.pattern}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	return conn, err
}

func (c *client) Send(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	// todo retry
	conn, err := c.dail()
	if err != nil {
		return err
	}
	c.Conn = conn

	return c.WriteMessage(websocket.TextMessage, data)
}

func (c *client) Read(v any) error {
	_, msg, err := c.Conn.ReadMessage()
	if err != nil {
		return err
	}

	return json.Unmarshal(msg, v)
}
