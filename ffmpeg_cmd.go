package vstreamer

import (
	"os"
	"io"
	"bytes"
	"os/exec"
	"log"
	"fmt"
)

// ffmpeg wrapper
type ffmpeg struct {
	// Video streamer server that will broadcast the video stream produced by ffmpeg
	listener *VServer
	stopMe chan bool
}

const(
	FFMPEG_CMD = "/usr/bin/ffmpeg"
)

var FFMPEG_ARGS = []string{"-i", "/dev/video2", "-f", "mpeg1video", "-"}

func (f *ffmpeg) Run(){
	buffer := make([]byte, 8096)
	ffmpeg := exec.Command(FFMPEG_CMD, FFMPEG_ARGS...)
    stdout, _ := ffmpeg.StdoutPipe()
	if err := ffmpeg.Start(); err != nil{
		log.Println(fmt.Sprintf("ffmpeg invocation failed CMD:%s ARGS:%s", FFMPEG_CMD, FFMPEG_ARGS))
		os.Exit(-1)
	}

	// Wait for ffmpeg to finish
	go func(){
		if err := ffmpeg.Wait(); err != nil{
			if ffmpeg.ProcessState.Exited() {
				log.Println("ffmpeg invocation failed")
				os.Exit(-1);
			} else {
				log.Println("ffmpeg was intentionally stopped")
			}
		}
	}()
	// Broadcast data
	go func(){
		for{
			count, err := stdout.Read(buffer)
			if err == io.EOF {
				stdout.Close()
				break
			}
			f.listener.Broadcast(bytes.NewReader(buffer[0:count]))
		}
	}()
	// Check if process should be terminated
	go func(){
		<-f.stopMe
		if err := ffmpeg.Process.Kill(); err != nil{
			log.Println("An error was encounter stopping ffmpeg process ", err)
		}
	}()
}

func (f *ffmpeg) Start(l *VServer){
	if f.listener == nil{
		f.listener = l
		log.Println("Starting stream ", *l.in)
		FFMPEG_ARGS[1] = *l.in
		f.Run()
	} else {
		log.Println("Stream already started")
	}
}

func (f *ffmpeg) Stop(){
	log.Println("Stopping stream")
	f.stopMe <- true
	f.listener = nil
}

func NewFfmpegProcess() *ffmpeg{
	ff := ffmpeg{
		listener: nil,
		stopMe: make(chan bool),
	}

	return &ff
}
