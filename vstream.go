package vstreamer

import (
	"bytes"
	"encoding/binary"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
)

const (
	socketBufferSize = 1024
)

type VServer struct {
	clients     map[*client]bool
	width       uint16
	height      uint16
	in          *string
	videoSource *ffmpeg
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

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
	client := NewClient(socket)

	s.clients[client] = true
	log.Println("Websocket clients:", len(s.clients))

	if len(s.clients) == 1 {
		s.videoSource.Start(s)
	}
}

func (s *VServer) Broadcast(reader *bytes.Reader) {
	var writers []io.Writer
	for client := range s.clients {
		writer, err := client.socket.NextWriter(websocket.BinaryMessage)
		if err != nil {
			delete(s.clients, client)
			log.Println("Websocket clients:", len(s.clients))
			if len(s.clients) == 0 {
				s.videoSource.Stop()
			}
			continue
		}
		writers = append(writers, writer)
		defer writer.Close()
	}

	if len(writers) > 0 {
		wrtr := io.MultiWriter(writers...)
		if _, err := io.Copy(wrtr, reader); err != nil {
			log.Println(err)
		}
	}
}

func NewServer(width uint16, height uint16, videoIn *string) *VServer {
	vs := VServer{
		clients:     make(map[*client]bool),
		width:       width,
		height:      height,
		in:          videoIn,
		videoSource: NewFfmpegProcess(),
	}
	return &vs
}
