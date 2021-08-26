// Package main to build CLI executable
package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/jusongchen/REST-app/pkg/exampleapp/cmd"
)

func main() {
	//needed for glog
	_ = flag.Set("logtostderr", "true")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	cmd.Execute()
}
