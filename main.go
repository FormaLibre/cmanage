package main

import (
	"github.com/forma-libre/cmanage/cmd"
)

var (
	// VERSION is set during build
	VERSION = "0.0.1"
)

func main() {
	cmd.Execute(VERSION)
}
