.PHONY: build run
default: build

build:
	GO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o release/benchmark-proxy_windows_amd64.exe main.go
	GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o release/benchmark-proxy_linux_amd64 main.go
	GO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o release/benchmark-proxy_mac_amd64 main.go

run:
	go run main.go