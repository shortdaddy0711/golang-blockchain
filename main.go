package main

import (
	"os"

	"github.com/shortdaddy0711/golang-blockchain/cli"
)

func main() {
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()
}