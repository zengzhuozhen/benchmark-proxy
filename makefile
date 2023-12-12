.PHONY: build run
default: build

build:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o release/benchmark-proxy_windows_amd64.exe main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o release/benchmark-proxy_linux_amd64 main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o release/benchmark-proxy_mac_amd64 main.go

test:
	go test -v ./...

run:
	go run main.go