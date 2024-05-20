package main

import (
	"os"

	"bitbucket.org/neiku/winornot"
	"github.com/PureMature/starcli/cli"
)

func init() {
	// fix for Windows terminal output
	winornot.EnableANSIControl()
}

func main() {
	args := cli.ParseArgs()
	os.Exit(cli.Process(args))
}
