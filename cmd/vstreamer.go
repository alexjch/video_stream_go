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
	staticWWW := path.Join(staticDir, "www")
    log.Println("Serving static files from: ", staticWWW)
    return http.FileServer(http.Dir(staticWWW))
}


func parseArgs() (int, *string, uint16, uint16){
	port := flag.Int("port", DEFAULT_PORT, "Port number")
	video_in := flag.String("video", DEFAULT_VIDEO_IN, "Video input")
	width := flag.Uint("width", DEFAULT_W, "Video width")
	height := flag.Uint("height", DEFAULT_H, "Video height")
    // TODO: serve width, height from vstreamer invocation
	flag.Parse()
	return *port, video_in, uint16(*width), uint16(*height)
}

func main(){
	port, video_in, width, height := parseArgs()
	// Validating if video input exists
	if _, err := os.Stat(*video_in); os.IsNotExist(err){
		log.Println(*video_in, " does not exists")
		os.Exit(-1)
	}
	server_addr := fmt.Sprintf("0.0.0.0:%d", port)
	addr := flag.String("addr", server_addr, "http service address")
	videoStreamer := vstreamer.NewServer(width, height, video_in)
	http.HandleFunc("/echo", videoStreamer.Echo)
    http.Handle("/", getStaticDir())
	log.Println("Starting web socket server on: ", server_addr)
	http.ListenAndServe(*addr, nil)
}
