#!/usr/bin/env bash

#go run -ldflags "-X main.Version=v1.0.1 -X 'main.BuildTime=$(date +'%Y/%m/%d %H:%M:%S')'" main.go

go build -ldflags "-X main.buildVersion=$1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')'" -o ./cmd/shortener ./cmd/shortener/.
