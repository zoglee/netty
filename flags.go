package main

import (
	"flag"
)

var (
	debug = false
	port  = "8080"
)

func init() {
	flag.BoolVar(&debug, "debug", false, "log debug info or not")
	flag.StringVar(&port, "port", "8080", "listen port")
}
