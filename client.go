package vstreamer

import (
	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
}

func (c *client) readLoop() {
	for {
		if _, _, err := c.socket.NextReader(); err != nil {
			c.socket.Close()
			break
		}
	}
}

func NewClient(socket *websocket.Conn) *client {
	c := client{
		socket: socket,
	}
	go c.readLoop()
	return &c
}
