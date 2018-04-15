package main

import "github.com/wjdp/htmltest/cmd"

var (
	version   string
	date string
)

func main() {
	cmd.Version = version
	cmd.BuildDate = date
	cmd.Execute()
}
