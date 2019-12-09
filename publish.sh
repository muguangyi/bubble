#!/bin/bash

set GOPROXY=https://goproxy.io
set CGO_ENABLED=0

echo --- Build for windows ---
set GOOS=windows
set GOARCH=amd64
go build -o pub/windows/master/bubble-master.exe bubble-master/main.go
go build -o pub/windows/worker/bubble-worker.exe bubble-worker/main.go
cp config/master.toml config/log.xml config/master.yml pub/windows/master/
cp config/worker.toml config/log.xml config/worker.yml pub/windows/worker/

echo --- Build for linux ---
set GOOS=linux
set GOARCH=amd64
go build -o pub/linux/master/bubble-master bubble-master/main.go
go build -o pub/linux/worker/bubble-worker bubble-worker/main.go
cp config/master.toml config/log.xml config/master.yml pub/linux/master/
cp config/worker.toml config/log.xml config/worker.yml pub/linux/worker/

echo --- Build for mac ---
set GOOS=darwin
set GOARCH=amd64
go build -o pub/mac/master/bubble-master bubble-master/main.go
go build -o pub/mac/worker/bubble-worker bubble-worker/main.go
cp config/master.toml config/log.xml config/master.yml pub/mac/master/
cp config/worker.toml config/log.xml config/worker.yml pub/mac/worker/

echo --- Build portal ---
cd bubble-portal
call npm install
call npm run build
cp -r ./dist ../pub/windows/master/
cp -r ./dist ../pub/linux/master/
cp -r ./dist ../pub/mac/master/
cd ..

echo === Publish completed ===
