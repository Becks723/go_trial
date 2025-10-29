# 编译阶段
# 镜像
FROM golang:1.25.1-alpine AS builder

# 工作目录
WORKDIR /app

# 设置国内 Go 模块代理
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 拷贝go.mod 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 拷贝源代码并构建
COPY . .
ENV CGO_ENABLED=0
RUN go build -o memo ./cmd

# 打包阶段
FROM alpine:latest
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /app/memo .

# 暴露服务端口
EXPOSE 8080

# 启动命令
CMD ["./memo"]
