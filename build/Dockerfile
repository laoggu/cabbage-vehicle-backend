#-------- 构建阶段 --------
FROM golang:1.21-alpine AS builder

# 安装 ca-certificates 供 HTTPS 调用（微信、OSS）
RUN apk add --no-cache ca-certificates git

WORKDIR /src

# 提前拷贝 go mod 并下载依赖，利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 再拷贝源码，避免每次依赖变动都全量重新编译
COPY . .
# 静态编译，减小体积并可在 scratch 运行
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags='-w -s' -o /bin/gateway ./cmd/gateway

#-------- 运行阶段 --------
FROM scratch
# 时区数据（可选，如需日志用本地时间）
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
# SSL 根证书
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# 可执行文件
COPY --from=builder /bin/gateway /gateway

ENTRYPOINT ["/gateway"]