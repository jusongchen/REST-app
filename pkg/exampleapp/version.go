package exampleapp

import (
	"fmt"
)

var (
	// BuildTime is a time label of the moment when the binary was built
	BuildTime = "unset"
	// Commit is a last commit hash at the moment when the binary was built
	Commit = "unset"
	// Release is a semantic version of current build
	Release = "unset"
)

//PrintBanner print version banner
func PrintBanner(decoration ...string) {
	fmt.Printf("Version: release %s,build time %s,commit %s\n", Release, BuildTime, Commit)
}
