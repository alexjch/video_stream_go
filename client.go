package vstreamer

import (
	"github.com/gorilla/websocket"
)

type client struct{
	socket *websocket.Conn
}
