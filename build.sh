#!/bin/bash

# Script for devs to build a copy of the app with build flags set

go build -ldflags "-X main.buildDate=`date -u +%Y-%m-%dT%H:%M:%SZ` -X main.version=`git describe --tags`" -o bin/htmltest -x main.go
