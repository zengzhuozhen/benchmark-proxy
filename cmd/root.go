package cmd

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/zengzhuozhen/benchmark-proxy/core"
	"github.com/zengzhuozhen/benchmark-proxy/resources"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	port    int
	isDebug bool
)

var rootCmd = &cobra.Command{
	Use:   "benchmark-proxy",
	Short: "Benchmark-Proxy is a proxy server for HTTP/HTTPS benchmark",
	Long: `Benchmark-Proxy is a proxy server for HTTP/HTTPS benchmark
                Use it by curl -x option like: curl -x 127.0.0.1:9900 https://www.baidu.com `,
	Run: func(cmd *cobra.Command, args []string) {
		if isDebug {
			log.SetLevel(log.DebugLevel)
			go pprof()
		}
		ca, key := parseCA()
		proxy := core.NewBenchProxyService(port, ca, key)
		go proxy.Serve()
		log.Infof("proxy started success in 127.0.0.1:%d \n", port)
		gracefulStop()
	},
}

func pprof() {
	if err := http.ListenAndServe("127.0.0.1:6060", nil); err != nil {
		log.Errorln("start pprof failed")
	}
}

func gracefulStop() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Infof("receive signal %s \n", sig)
	log.Infoln("Graceful Exit")
	os.Exit(0)
}

func parseCA() (rootCA *x509.Certificate, rootKey *rsa.PrivateKey) {
	var err error
	block, _ := pem.Decode(resources.CaCrt)
	rootCA, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(fmt.Errorf("加载根证书失败: %s", err))
	}
	block, _ = pem.Decode(resources.CaKey)
	rootKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(fmt.Errorf("加载根证书私钥失败: %s", err))
	}
	log.Debugln("加载证书成功")
	return
}

func init() {
	rootCmd.PersistentFlags().IntVar(&port, "port", 9900, "proxy server bind port")
	rootCmd.PersistentFlags().BoolVar(&isDebug, "debug", false, "debug mode")

	// subCmd
	rootCmd.AddCommand(tagsCmd, headersCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}
