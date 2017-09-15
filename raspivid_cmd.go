// +build linux
// +build arm

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

// raspivid wrapper
type VideoStream struct {
	// Video streamer server that will broadcast the video stream produced by raspivid
	listener *VServer
	stopMe   chan bool
}

const (
	RASPIVID_CMD = "/usr/bin/raspivid"
)

var RASPIVID_ARGS = []string{"--nopreview", "--exposure", "auto", "--output", "-"}

func (f *VideoStream) Run() {
    var wg sync.WaitGroup
	buffer := make([]byte, 8096)
	raspivid := exec.Command(RASPIVID_CMD, RASPIVID_ARGS...)
	stdout, _ := raspivid.StdoutPipe()
	if err := raspivid.Start(); err != nil {
		log.Println(fmt.Sprintf("raspivid invocation failed CMD:%s ARGS:%s", RASPIVID_CMD, RASPIVID_ARGS))
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
        log.Println("Sending SIGINT to raspivid")
        if err := raspivid.Process.Signal(syscall.SIGINT); err != nil {
			log.Println("An error was encounter stopping raspivid process ", err)
		}
        wg.Done()
	}()
    // Wait for raspivid to finish
    raspivid.Wait()
    wg.Wait()
    f.listener = nil
}

func (f *VideoStream) Start(l *VServer) {
	if f.listener == nil {
		f.listener = l
		log.Println("Starting stream ", *l.in)
		RASPIVID_ARGS[1] = *l.in
		f.Run()
	} else {
		log.Println("Stream already started")
	}
}

func (f *VideoStream) Stop() {
	log.Println("Stopping stream")
	f.stopMe <- true
}

func NewVideoSource() *VideoStream {
	ff := VideoStream{
		listener: nil,
		stopMe:   make(chan bool),
	}

	return &ff
}
