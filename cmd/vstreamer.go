package main

import (
	"fmt"
	"log"
	"flag"
	"net/http"
	"github.com/alexjch/vstreamer"
	"strconv"
)

const (
	DEFAULT_PORT = "7072"
	DEFAULT_VIDEO_IN = "/dev/video0"
	DEFAULT_W = "640"
	DEFAULT_H = "480"
	usage = "Usage: vstreamer [W=<width> H=<height> P=<port> V=<video_in>]"
)

func parseArgs() (*string, *string, uint16, uint16){
	port := flag.String("P", DEFAULT_PORT, "Port number")
	video_in := flag.String("video", DEFAULT_VIDEO_IN, "Video input")
	width := flag.String("W", DEFAULT_W, "Video width")
	height := flag.String("H", DEFAULT_H, "Video height")
	w, _ := strconv.ParseUint(*width, 10, 16)
	h, _ := strconv.ParseUint(*height, 10, 16)
	flag.Parse()
	return port, video_in, uint16(w), uint16(h)
}

func main(){
	port, video_in, width, height := parseArgs()
	server_addr := fmt.Sprintf("0.0.0.0:%s", *port)
	addr := flag.String("addr", server_addr, "http service address")
	videoStreamer := vstreamer.NewServer(width, height)
	videoSource := vstreamer.NewFfmpegProcess(videoStreamer)
	videoSource.Start(video_in)
	http.HandleFunc("/echo", videoStreamer.Echo)
	log.Println("Starting web socket server on: ", server_addr)
	http.ListenAndServe(*addr, nil)
}
