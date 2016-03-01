package main

import (
	"log"
	"flag"
	"net/http"
	"fmt"
	"bytes"
	"github.com/alexjch/vstreamer"
)

const(
	FFMPEG_CMD = "/usr/bin/ffmpeg"
)
//var FFMPEG_ARGS = []string{"-i", "/dev/video0", "-f", "mpeg1video", "-"}
var FFMPEG_ARGS = []string{"-i", "BigBuckBunny_640x360.m4v", "-f", "mpeg1video", "-"}
//var addr = flag.String("addr", "10.0.1.21:8080", "http service address")
var addr = flag.String("addr", "0.0.0.0:3000", "http service address")

func setSource(s *vstreamer.VServer) http.HandlerFunc{
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		p := make([]byte, 8096)
		for{
			count, err := r.Body.Read(p);
			if err != nil{
				log.Println(err)
				r.Body.Close()
				return
			}
			s.Broadcast(bytes.NewReader(p[0:count]))
		}
	})
}

func main(){
	fmt.Println("Starting streamming")
	var videoStreamer = vstreamer.NewServer()
	http.HandleFunc("/echo", videoStreamer.Echo)
	http.HandleFunc("/stream", setSource(videoStreamer))
	http.ListenAndServe(*addr, nil)
}
