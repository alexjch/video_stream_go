package vstreamer

import (
	"io"
	"log"
	"bytes"
	"net/http"
	"encoding/binary"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

type VServer struct{
	clients map[*client]bool
	width uint16
	height uint16
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

/*
var width uint16 = 640
var height uint16 = 480
*/
var magicBytes = []byte("jsmp")

func (s *VServer) Echo(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, magicBytes)
	binary.Write(buf, binary.BigEndian, s.width)
	binary.Write(buf, binary.BigEndian, s.height)

	err = socket.WriteMessage(websocket.BinaryMessage, buf.Bytes())

	client := &client{
		socket: socket,
	}

	s.clients[client] = true
	log.Println("Websocket clients:", len(s.clients))
}

func (s *VServer) Broadcast(reader *bytes.Reader) {
	for client := range s.clients{
		writer, _ := client.socket.NextWriter(websocket.BinaryMessage)
		if _, err := io.Copy(writer, reader); err != nil{
			delete(s.clients, client)
		}
		writer.Close()
	}
}

func NewServer(width uint16, height uint16) *VServer{
	vs := VServer{
		clients: make(map[*client]bool),
		width: width,
		height: height,
	}
	return &vs
}