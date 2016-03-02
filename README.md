# video_stream_go
Streaming video to websockets using go, inspired in [this project](http://phoboslab.org/log/2013/09/html5-live-video-streaming-via-websockets)

This branch waits for a stream of video to connect to the server on <host>:3000/stream, the stream can be started using: 

ffmpeg -i /vagrant/BigBuckBunny_640x360.m4v -f mpeg1video http://localhost:3000/stream
