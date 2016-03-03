# video_stream_go
Streaming video to websockets using go, inspired by [this project](http://phoboslab.org/log/2013/09/html5-live-video-streaming-via-websockets). Modifications
were made to the original idea and a whole new server using go was created to my specific usage (IP camera using intel edison)

The default options will stream a 640x480 stream from /dev/video0

```
# > ./vstreamer
```

Possible command line arguments are:

```
-video=<video source> -heigh=<video_height_format> -width=<video_width_format> -port=<server_port>
```

Where "Video Source" could be /dev/video0 or /tmp/BigBuckBunny_640x360.m4v, "Width and Height" are the dimensions of the video frame, and "Server Port" the port
to use to serve static files and the stream using websockets.

Once the server is compiled make sure you have a folder structure similar to the following:

```
root/
    vstreamer
    www/
       index.html
       js/
         jsmpg.js
```

The binary vstreamer will look for the static files (`.html` and `.js`) under `www` and `www/js` folders respectively.


