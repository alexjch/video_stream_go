package main

import (
    "os"
	"fmt"
	"log"
	"flag"
    "path"
    "path/filepath"
	"net/http"
	"github.com/alexjch/vstreamer"
)

const (
	DEFAULT_PORT = 7072
	DEFAULT_VIDEO_IN = "/dev/video0"
	DEFAULT_W = 640
	DEFAULT_H = 480
	usage = "Usage: vstreamer [-width=<width> -height=<height> -port=<port> -video=<video_in>]"
)

func getStaticDir() http.Handler{
    staticDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
    log.Println(staticDir)
    return http.FileServer(http.Dir(path.Join(staticDir, "www")))
}


func parseArgs() (int, *string, uint16, uint16){
	port := flag.Int("port", DEFAULT_PORT, "Port number")
	video_in := flag.String("video", DEFAULT_VIDEO_IN, "Video input")
	width := flag.Uint("width", DEFAULT_W, "Video width")
	height := flag.Uint("height", DEFAULT_H, "Video height")
    // TODO: handle error and print usage
    // TODO: serve static content
    // TODO: serve width, height from movie
    // TODO: shutdown, startup video feed when clients = 0
	flag.Parse()
	return *port, video_in, uint16(*width), uint16(*height)
}

func main(){
	port, video_in, width, height := parseArgs()
	server_addr := fmt.Sprintf("0.0.0.0:%d", port)
	addr := flag.String("addr", server_addr, "http service address")
	videoStreamer := vstreamer.NewServer(width, height)
	videoSource := vstreamer.NewFfmpegProcess(videoStreamer)
	videoSource.Start(video_in)
	http.HandleFunc("/echo", videoStreamer.Echo)
    http.Handle("/", getStaticDir())
	log.Println("Starting web socket server on: ", server_addr)
	http.ListenAndServe(*addr, nil)
}
