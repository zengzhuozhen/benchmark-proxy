# 使用官方 Go 镜像作为基础镜像
FROM golang:1.18-alpine

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o benchmark-proxy .

# 暴露端口（根据您的应用需要修改端口号）
EXPOSE 9900

# 运行应用
CMD ["./benchmark-proxy"]