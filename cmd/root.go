package cmd

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zengzhuozhen/benchmark-proxy/core"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	port    int
	rootCA  string
	rootKey string
	isDebug bool
)

var rootCmd = &cobra.Command{
	Use:   "benchmark-proxy",
	Short: "Benchmark-Proxy is a proxy server for HTTP/HTTPS benchmark",
	Long: `Benchmark-Proxy is a proxy server for HTTP/HTTPS benchmark
                Use it by curl -x option like: curl -x 127.0.0.1:9900 https://www.baidu.com `,
	Run: func(cmd *cobra.Command, args []string) {
		ca, key := parseCA(rootCA, rootKey)
		proxy := core.NewBenchProxyService(port, ca, key)
		go proxy.Serve(isDebug)
		fmt.Printf("proxy started success in 127.0.0.1:%d \n", port)
		if isDebug {
			go pprof()
			fmt.Println("「DEBUG」Begin...")
		}
		gracefulStop()
	},
}

func pprof() {
	if err := http.ListenAndServe("127.0.0.1:6060", nil); err != nil {
		fmt.Printf("start pprof failed")
	}
}

func gracefulStop() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	fmt.Printf("receive signal %s \n", sig)
	fmt.Println("Graceful Exit")
	os.Exit(0)
}

func checkFilesExist(files ...string) bool {
	for _, i := range files {
		_, err := os.Stat(i)
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func parseCA(crt, key string) (rootCA *x509.Certificate, rootKey *rsa.PrivateKey) {
	var err error
	if !checkFilesExist(crt, key) {
		return
	}
	crtByte, err := ioutil.ReadFile(crt)
	if err != nil {
		panic(fmt.Errorf("加载证书文件{ca-crt}失败: %s", err))
	}
	keyByte, err := ioutil.ReadFile(key)
	if err != nil {
		panic(fmt.Errorf("加载证书文件{ca-key}失败: %s", err))
	}
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
	if isDebug {
		fmt.Println("加载证书成功")
	}
	return
}

func init() {
	rootCmd.PersistentFlags().IntVar(&port, "port", 9900, "proxy server bind port")
	rootCmd.PersistentFlags().StringVar(&rootCA, "ca-crt", "ca.crt", "ca.crt file for HTTPS proxy,default: 'ca.crt' in root dir")
	rootCmd.PersistentFlags().StringVar(&rootKey, "ca-key", "ca.key", "ca.crt file for HTTPS proxy,default: 'ca.key' in root dir")
	rootCmd.MarkFlagsRequiredTogether("ca-crt", "ca-key")
	rootCmd.PersistentFlags().BoolVar(&isDebug, "debug", false, "debug mode")

	// subCmd
	rootCmd.AddCommand(tagsCmd, headersCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
