package vstreamer

import (
	"os"
	"io"
	"bytes"
	"os/exec"
	"log"
)

// ffmpeg wrapper
type ffmpeg struct {
	// Video streamer server that will broadcast the video stream produced by ffmpeg
	listener *VServer
}

const(
	FFMPEG_CMD = "/usr/bin/ffmpeg"
)

var FFMPEG_ARGS = []string{"-i", "/dev/video0", "-f", "mpeg1video", "-"}

func (f *ffmpeg) Run(){
	buffer := make([]byte, 8096)
	ffmpeg := exec.Command(FFMPEG_CMD, FFMPEG_ARGS...)
    stdout, _ := ffmpeg.StdoutPipe()
	ffmpeg.Start()

	go func(){
		if err := ffmpeg.Wait(); err != nil{
			log.Println("ffmpeg invocation failed, make sure the input video exists")
			os.Exit(-1);
		}
	}()

	for{
		count, err := stdout.Read(buffer)
		if err == io.EOF {
			stdout.Close()
			break
		}
		f.listener.Broadcast(bytes.NewReader(buffer[0:count]))
	}
}

func (f *ffmpeg) Start(in *string){
	log.Println("Starting stream ", *in)
	FFMPEG_ARGS[1] = *in
	go f.Run()
}

func (f *ffmpeg) Stop(){

}

func NewFfmpegProcess(s *VServer) *ffmpeg{
	ff := ffmpeg{
		listener: s,
	}

	return &ff
}
