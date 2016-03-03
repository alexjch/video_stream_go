package vstreamer

import (
	"io"
	"bytes"
	"os/exec"
	"log"
)

type ffmpeg struct {
	listener *VServer
}

const(
	FFMPEG_CMD = "/usr/bin/ffmpeg"
)

var FFMPEG_ARGS = []string{"-i", "/dev/video0", "-f", "mpeg1video", "-"}

func ffmpegStart(s *VServer){
	buffer := make([]byte, 8096)
	ffmpeg := exec.Command(FFMPEG_CMD, FFMPEG_ARGS...)
    stdout, _ := ffmpeg.StdoutPipe()
	_ = ffmpeg.Start()

	for{
		count, err := stdout.Read(buffer)
		if err == io.EOF {
			stdout.Close()
			break
		}
		s.Broadcast(bytes.NewReader(buffer[0:count]))
	}
}

func (f *ffmpeg) Start(in *string){
	log.Println("Starting stream with: ", *in)
	FFMPEG_ARGS[1] = *in
	go ffmpegStart(f.listener)
}

func (f *ffmpeg) Stop(){

}

func NewFfmpegProcess(s *VServer) *ffmpeg{
	ff := ffmpeg{
		listener: s,
	}

	return &ff
}
