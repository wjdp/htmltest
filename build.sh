#!/bin/bash

go build -ldflags "-X main.buildDate=`date -u +%Y-%m-%dT%H:%M:%SZ`" -o bin/htmltest -x main.go
