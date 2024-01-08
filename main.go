package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/zengzhuozhen/benchmark-proxy/cmd"
	_ "net/http/pprof"
	"os"
)

func main() {
	cmd.Execute()
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
