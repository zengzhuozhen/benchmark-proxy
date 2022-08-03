package main

import (
	"github.com/zengzhuozhen/benchmark-proxy/core"
	"flag"
	"fmt"
)

const helpDesc = `
   http benchmark proxy
`

func main() {
	flag.Usage = func() {
		fmt.Println(helpDesc)
		flag.PrintDefaults()
	}
	port := flag.Int("p", 9900, "proxy port")
	proxy := core.NewBenchProxyService(*port)
	proxy.Serve()
}
