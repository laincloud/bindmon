#!/bin/bash

docker run --rm -v $GOPATH:/go -e GOBIN=/go/src/github.com/laincloud/bindmon/bin golang:1.10 go install github.com/laincloud/bindmon