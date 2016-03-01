package vstreamer

import (
	"github.com/gorilla/websocket"
)

/**
reverse port forwarding in vagrant
vagrant ssh -- -R 12345:localhost:8082
ffmpeg -i /vagrant/BigBuckBunny_640x360.m4v -f mpeg1video http://localhost:12345/test/640/480
 */

type client struct{
	socket *websocket.Conn
}
