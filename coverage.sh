#!/bin/sh

go test ./htmldoc -coverprofile htmldoc.out
go test ./htmltest -coverprofile htmltest.out
go test ./issues -coverprofile issues.out
go test ./refcache -coverprofile refcache.out

go tool cover -html htmldoc.out
go tool cover -html htmltest.out
go tool cover -html issues.out
go tool cover -html refcache.out
