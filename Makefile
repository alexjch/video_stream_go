
edison:
	GOOS=linux GOARCH=386 go build cmd/vstreamer.go

pi:
	GOOS=linux GOARCH=arm GOARM=6 go build cmd/vstreamer.go


