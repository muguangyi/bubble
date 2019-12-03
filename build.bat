@echo off

set GOPROXY=https://goproxy.io
set CGO_ENABLED=0

echo --- Build for windows ---
set GOOS=windows
set GOARCH=amd64
go build -o bin/master/bubble-master.exe bubble-master/main.go
go build -o bin/worker/bubble-worker.exe bubble-worker/main.go

