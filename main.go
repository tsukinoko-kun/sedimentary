package main

import (
	"github.com/tsukinoko-kun/sedimentary/build"
	"github.com/tsukinoko-kun/sedimentary/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	build.Version = version
	build.Commit = commit
	build.Date = date

	cmd.Execute()
}
