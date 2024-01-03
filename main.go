package main

import "github.com/zengzhuozhen/benchmark-proxy/cmd"
import _ "net/http/pprof"

// todo 透明http代理，启动写入环境变量
// todo gomebd写入ca
func main() {
	cmd.Execute()
}
