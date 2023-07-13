package main

import "github.com/zengzhuozhen/benchmark-proxy/cmd"
import _ "net/http/pprof"

func main() {
	cmd.Execute()
}
