########## 变量 ##########
BIN_NAME      := gateway
DOCKER_COMPOSE:= docker-compose -f build/docker-compose.dev.yml
BUF           := buf
GQLGEN        := go run github.com/99designs/gqlgen generate
BENCH_PATH    := ./test/bench/gateway

########## 默认目标 ##########
.PHONY: help
help:  ## 显示可用命令
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

########## 代码生成 ##########
.PHONY: proto
proto: ## 生成 pb 代码（依赖 buf）
	@$(BUF) generate

.PHONY: gqlgen
gqlgen: ## 生成 GraphQL 代码
	@$(GQLGEN)

########## 编译 ##########
.PHONY: build
build: ## 编译 gateway 二进制
	@go build -ldflags="-w -s" -o bin/$(BIN_NAME) ./cmd/gateway

########## 本地环境 ##########
.PHONY: up
up: ## 启动依赖（MySQL/Redis/Kafka）
	@$(DOCKER_COMPOSE) up -d

.PHONY: down
down: ## 停止依赖
	@$(DOCKER_COMPOSE) down

.PHONY: logs
logs: ## 查看依赖日志
	@$(DOCKER_COMPOSE) logs -f

########## 运行 ##########
.PHONY: run
run: build up ## 编译+启动依赖+跑本地网关（前台）
	./bin/$(BIN_NAME)

########## 测试 ##########
.PHONY: bench
bench: ## 压测网关（需提前起服务）
	@mkdir -p $(BENCH_PATH)
	go test -bench=. -benchmem $(BENCH_PATH)

.PHONY: test
test: ## 跑单元测试
	go test ./...

########## 工具 ##########
.PHONY: fmt
fmt: ## 格式化代码
	@gofmt -s -w .

.PHONY: lint
lint: ## 静态检查（需安装 golangci-lint）
	@golangci-lint run

.PHONY: clean
clean: ## 清理二进制和缓存
	@rm -rf bin/
	@go clean -cache -testcache