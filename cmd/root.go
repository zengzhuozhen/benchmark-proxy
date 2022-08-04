package cmd

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zengzhuozhen/benchmark-proxy/core"
	"io/ioutil"
	"os"
)

var (
	port    int
	rootCA  string
	rootKey string
)

var rootCmd = &cobra.Command{
	Use:   "benchmark-proxy",
	Short: "benchmark-proxy is a proxy server for HTTP/HTTPS benchmark",
	Long: `benchmark-proxy is a proxy server for HTTP/HTTPS benchmark
                use it by curl -x option like:
                curl -x 127.0.0.1:9900 http://www.baidu.com`,
	Run: func(cmd *cobra.Command, args []string) {
		ca, key := parseCA(rootCA, rootKey)
		proxy := core.NewBenchProxyService(port, ca, key)
		proxy.Serve()
	},
}

func parseCA(crt, key string) (rootCA *x509.Certificate, rootKey *rsa.PrivateKey) {
	var err error
	crtByte, _ := ioutil.ReadFile(crt)
	keyByte, _ := ioutil.ReadFile(key)
	block, _ := pem.Decode(crtByte)
	rootCA, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(fmt.Errorf("加载根证书失败: %s", err))
	}
	block, _ = pem.Decode(keyByte)
	rootKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(fmt.Errorf("加载根证书私钥失败: %s", err))
	}
	return
}

func init() {
	rootCmd.PersistentFlags().IntVar(&port, "port", 9900, "proxy server bind port")
	rootCmd.PersistentFlags().StringVar(&rootCA, "ca-crt", "ca.crt", "ca.crt file for HTTPS proxy,default: 'ca.crt' is root dir")
	rootCmd.PersistentFlags().StringVar(&rootKey, "ca-key", "ca.key", "ca.crt file for HTTPS proxy,default: 'ca.key' is root dir")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
