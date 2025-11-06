# Makefile for Kafka Demo Project

.PHONY: help setup start stop clean install test producer consumer

# 默认目标
help:
	@echo "Kafka + Go 学习项目"
	@echo ""
	@echo "可用命令:"
	@echo "  make setup       - 初始化项目（下载依赖）"
	@echo "  make start       - 启动 Kafka 集群"
	@echo "  make stop        - 停止 Kafka 集群"
	@echo "  make restart     - 重启 Kafka 集群"
	@echo "  make clean       - 清理所有数据"
	@echo "  make logs        - 查看 Kafka 日志"
	@echo "  make status      - 查看集群状态"
	@echo ""
	@echo "示例运行:"
	@echo "  make producer    - 运行简单生产者"
	@echo "  make consumer    - 运行简单消费者"
	@echo "  make async       - 运行异步生产者"
	@echo "  make batch       - 运行批量消费者"
	@echo "  make order       - 运行订单处理系统"
	@echo ""
	@echo "Kafka 管理:"
	@echo "  make topics      - 列出所有 Topic"
	@echo "  make groups      - 列出所有消费者组"
	@echo "  make ui          - 打开 Kafka UI"

# 初始化项目
setup:
	@echo "📦 下载 Go 依赖..."
	go mod download
	go mod tidy
	@echo "✅ 项目初始化完成"

# 启动 Kafka 集群
start:
	@echo "🚀 启动 Kafka 集群..."
	docker-compose up -d
	@echo "⏳ 等待服务启动（30秒）..."
	@sleep 30
	@echo "✅ Kafka 集群已启动"
	@echo "   Kafka UI: http://localhost:8080"
	@echo "   Kafdrop:  http://localhost:9000"

# 停止 Kafka 集群
stop:
	@echo "⏸️  停止 Kafka 集群..."
	docker-compose stop
	@echo "✅ Kafka 集群已停止"

# 重启 Kafka 集群
restart: stop start

# 清理数据和容器
clean:
	@echo "🧹 清理 Kafka 集群和数据..."
	docker-compose down -v
	@echo "✅ 清理完成"

# 查看日志
logs:
	docker-compose logs -f

# 查看集群状态
status:
	@echo "📊 集群状态:"
	@docker-compose ps

# 列出所有 Topic
topics:
	@echo "📋 Topic 列表:"
	@docker exec kafka1 kafka-topics.sh --list --bootstrap-server localhost:9092

# 列出所有消费者组
groups:
	@echo "👥 消费者组列表:"
	@docker exec kafka1 kafka-consumer-groups.sh --list --bootstrap-server localhost:9092

# 打开 Kafka UI
ui:
	@echo "🌐 打开 Kafka UI..."
	@open http://localhost:8080 2>/dev/null || xdg-open http://localhost:8080 2>/dev/null || echo "请手动打开: http://localhost:8080"

# 运行简单生产者
producer:
	@echo "🔵 运行简单生产者..."
	go run examples/01-simple-producer/main.go

# 运行简单消费者
consumer:
	@echo "🟢 运行简单消费者..."
	go run examples/02-simple-consumer/main.go

# 运行异步生产者
async:
	@echo "🔵 运行异步生产者..."
	go run examples/04-async-producer/main.go

# 运行批量消费者
batch:
	@echo "🟢 运行批量消费者..."
	go run examples/05-batch-consumer/main.go

# 运行订单处理系统
order:
	@echo "📦 运行订单处理系统..."
	go run examples/08-order-processing/main.go

# 测试所有代码
test:
	@echo "🧪 运行测试..."
	go test -v ./...

# 构建所有示例
build:
	@echo "🔨 构建所有示例..."
	@mkdir -p bin
	go build -o bin/simple-producer examples/01-simple-producer/main.go
	go build -o bin/simple-consumer examples/02-simple-consumer/main.go
	go build -o bin/async-producer examples/04-async-producer/main.go
	go build -o bin/batch-consumer examples/05-batch-consumer/main.go
	go build -o bin/order-service examples/08-order-processing/main.go
	@echo "✅ 构建完成，二进制文件在 bin/ 目录"

# 安装开发工具
install-tools:
	@echo "🔧 安装开发工具..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ 工具安装完成"

# 代码检查
lint:
	@echo "🔍 代码检查..."
	golangci-lint run ./...

# 格式化代码
fmt:
	@echo "💅 格式化代码..."
	go fmt ./...
	@echo "✅ 代码格式化完成"
