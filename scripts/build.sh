#!/bin/bash
set -e
[ -z "$1" ] && {
    echo must use with tag name
    exit 1
}
export GO111MODULE=on
export GOPROXY=GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,https://goproxy.io,https://athens.azurefd.net,direct
export CGO_ENABLED=0
go build -o dp main.go
GOOS=windows go build -o dp.exe main.go
tar zcf dp-windows-amd64-${1}.tar.gz dp.exe
tar zcf dp-linux-amd64-${1}.tar.gz dp

