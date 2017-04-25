package vstreamer

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
    "syscall"
    "sync"
)

// ffmpeg wrapper
type ffmpeg struct {
	// Video streamer server that will broadcast the video stream produced by ffmpeg
	listener *VServer
	stopMe   chan bool
}

const (
	FFMPEG_CMD = "/usr/bin/ffmpeg"
)

var FFMPEG_ARGS = []string{"-i", "/dev/video2", "-f", "mpeg1video", "-"}

func (f *ffmpeg) Run() {
    var wg sync.WaitGroup
	buffer := make([]byte, 8096)
	ffmpeg := exec.Command(FFMPEG_CMD, FFMPEG_ARGS...)
	stdout, _ := ffmpeg.StdoutPipe()
	if err := ffmpeg.Start(); err != nil {
		log.Println(fmt.Sprintf("ffmpeg invocation failed CMD:%s ARGS:%s", FFMPEG_CMD, FFMPEG_ARGS))
		os.Exit(-1)
	}
    wg.Add(2)
	// Broadcast data
	go func() {
		for {
			count, err := stdout.Read(buffer)
			if err == io.EOF {
				stdout.Close()
				break
			}
			f.listener.Broadcast(bytes.NewReader(buffer[0:count]))
		}
        wg.Done()
	}()
	// Check if process should be terminated
	go func() {
		<-f.stopMe
        log.Println("Sending SIGINT to ffmpeg")
        if err := ffmpeg.Process.Signal(syscall.SIGINT); err != nil {
			log.Println("An error was encounter stopping ffmpeg process ", err)
		}
        wg.Done()
	}()
    // Wait for ffmpeg to finish
    ffmpeg.Wait()
    wg.Wait()
    f.listener = nil
}

func (f *ffmpeg) Start(l *VServer) {
	if f.listener == nil {
		f.listener = l
		log.Println("Starting stream ", *l.in)
		FFMPEG_ARGS[1] = *l.in
		f.Run()
	} else {
		log.Println("Stream already started")
	}
}

func (f *ffmpeg) Stop() {
	log.Println("Stopping stream")
	f.stopMe <- true
}

func NewFfmpegProcess() *ffmpeg {
	ff := ffmpeg{
		listener: nil,
		stopMe:   make(chan bool),
	}

	return &ff
}
