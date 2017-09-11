#!/usr/bin/env bash

set -e
echo "mode: atomic" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -v -race -coverprofile=profile.out -covermode=atomic -coverpkg=./... $d
    if [ -f profile.out ]; then
        grep -h -v "^mode:" profile.out >> coverage.txt
        rm profile.out
    fi
done

# Print out the coverage.txt file with a total value
go tool cover -func=coverage.txt
