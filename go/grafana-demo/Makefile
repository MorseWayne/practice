.PHONY: help up down restart logs build run test clean

help: ## 显示帮助信息
	@echo "可用命令："
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

up: ## 启动所有服务
	docker-compose up -d
	@echo "✅ 服务已启动"
	@echo "📱 Go 应用: http://localhost:8080"
	@echo "📊 Prometheus: http://localhost:9090"
	@echo "📈 Grafana: http://localhost:3000 (admin/admin)"

down: ## 停止所有服务
	docker-compose down
	@echo "✅ 服务已停止"

restart: down up ## 重启所有服务

logs: ## 查看服务日志
	docker-compose logs -f

build: ## 构建 Go 应用
	go build -o bin/grafana-demo main.go
	@echo "✅ 构建完成: bin/grafana-demo"

run: ## 本地运行 Go 应用
	go run main.go

test: ## 运行测试
	go test -v ./...

clean: ## 清理容器和数据卷
	docker-compose down -v
	rm -rf bin/
	@echo "✅ 清理完成"

rebuild: ## 重新构建并启动
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d
	@echo "✅ 重新构建并启动完成"

ps: ## 查看运行状态
	docker-compose ps

deps: ## 下载 Go 依赖
	go mod download
	go mod tidy

fmt: ## 格式化代码
	go fmt ./...

vet: ## 代码检查
	go vet ./...
