# 快速入门指南

欢迎使用 Kafka + Go 学习项目！本指南将帮助您快速上手。

## 📋 前置要求

- **Go 1.21+**: [安装 Go](https://golang.org/dl/)
- **Docker & Docker Compose**: [安装 Docker](https://docs.docker.com/get-docker/)
- **至少 4GB 可用内存**

验证安装：
```bash
go version
docker --version
docker-compose --version
```

## 🚀 5 分钟快速开始

### 第 1 步：克隆项目（如果需要）

```bash
cd /home/quzhihao/workspace/source/private/practice/go/kafka-demo
```

### 第 2 步：初始化项目

```bash
make setup
```

这会下载所有 Go 依赖包。

### 第 3 步：启动 Kafka 集群

```bash
make start
```

等待约 30 秒，服务将完全启动。

### 第 4 步：运行第一个示例

**终端 1** - 启动生产者：
```bash
make producer
```

**终端 2** - 启动消费者：
```bash
make consumer
```

您将看到生产者发送消息，消费者接收消息的实时日志！

### 第 5 步：查看 Web UI

打开浏览器访问：
- **Kafka UI**: http://localhost:8080
- **Kafdrop**: http://localhost:9000

在这里可以可视化查看 Topic、消息、消费者组等信息。

## 📚 学习路径

### 第 1 天：基础概念

1. 阅读文档：
   - [Kafka 核心概念](./docs/01-kafka-basics.md)
   - [Kafka 架构设计](./docs/02-kafka-architecture.md)

2. 运行基础示例：
   ```bash
   # 生产者
   make producer
   
   # 消费者
   make consumer
   ```

### 第 2 天：深入理解

1. 阅读：
   - [消息模型与分区](./docs/03-message-partition.md)

2. 实验分区策略：
   - 修改 `examples/01-simple-producer/main.go`
   - 尝试不同的 Key 值
   - 观察消息分布

3. 测试消费者组：
   ```bash
   # 启动 3 个消费者实例
   # 终端 1
   make consumer
   
   # 终端 2
   make consumer
   
   # 终端 3
   make consumer
   ```
   观察分区是如何分配的。

### 第 3 天：高级特性

1. 异步生产者：
   ```bash
   make async
   ```
   学习批量发送和性能优化。

2. 批量消费者：
   ```bash
   # 先启动生产者
   make async
   
   # 再启动批量消费者
   make batch
   ```

### 第 4-5 天：实战项目

运行完整的订单处理系统：

```bash
# 终端 1: 订单服务
make order

# 观察整个业务流程
```

## 🎯 常用命令

### 项目管理
```bash
make setup      # 初始化项目
make build      # 构建所有示例
make test       # 运行测试
make fmt        # 格式化代码
```

### Kafka 集群
```bash
make start      # 启动集群
make stop       # 停止集群
make restart    # 重启集群
make clean      # 清理数据
make status     # 查看状态
make logs       # 查看日志
```

### 运行示例
```bash
make producer   # 简单生产者
make consumer   # 简单消费者
make async      # 异步生产者
make batch      # 批量消费者
make order      # 订单系统
```

### Kafka 管理
```bash
make topics     # 列出所有 Topic
make groups     # 列出消费者组
make ui         # 打开 Web UI
```

## 📖 示例目录

| 示例 | 说明 | 难度 |
|------|------|------|
| [01-simple-producer](./examples/01-simple-producer/) | 简单生产者 | ⭐ |
| [02-simple-consumer](./examples/02-simple-consumer/) | 简单消费者 | ⭐ |
| [04-async-producer](./examples/04-async-producer/) | 异步生产者 | ⭐⭐ |
| [05-batch-consumer](./examples/05-batch-consumer/) | 批量消费者 | ⭐⭐ |
| [08-order-processing](./examples/08-order-processing/) | 订单处理系统 | ⭐⭐⭐ |

## 🔧 Kafka 命令行工具

### 创建 Topic
```bash
docker exec kafka1 kafka-topics.sh --create \
  --bootstrap-server localhost:9092 \
  --topic my-topic \
  --partitions 3 \
  --replication-factor 2
```

### 列出 Topic
```bash
docker exec kafka1 kafka-topics.sh --list \
  --bootstrap-server localhost:9092
```

### 查看 Topic 详情
```bash
docker exec kafka1 kafka-topics.sh --describe \
  --bootstrap-server localhost:9092 \
  --topic example-topic
```

### 发送测试消息
```bash
docker exec -it kafka1 kafka-console-producer.sh \
  --bootstrap-server localhost:9092 \
  --topic example-topic
```

### 消费消息
```bash
docker exec -it kafka1 kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic example-topic \
  --from-beginning
```

### 查看消费者组
```bash
docker exec kafka1 kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --group example-consumer-group \
  --describe
```

## 🐛 故障排查

### 问题：无法连接到 Kafka

**解决方法：**
```bash
# 检查服务状态
make status

# 查看日志
make logs

# 重启集群
make restart
```

### 问题：端口已被占用

**解决方法：**
```bash
# 查看占用端口的进程
lsof -i :19092
lsof -i :8080

# 停止冲突的服务或修改 docker-compose.yml 中的端口
```

### 问题：消费者无法接收消息

**检查清单：**
1. ✅ Kafka 集群是否正常运行：`make status`
2. ✅ Topic 是否存在：`make topics`
3. ✅ 生产者是否成功发送：查看生产者日志
4. ✅ 消费者组配置是否正确
5. ✅ Offset 位置是否正确（尝试从头消费）

### 问题：Docker 内存不足

**解决方法：**
```bash
# 增加 Docker 内存限制（Docker Desktop）
# Settings -> Resources -> Memory: 至少 4GB

# 或减少 Kafka Broker 数量
# 编辑 docker-compose.yml，只保留 kafka1
```

## 💡 学习建议

1. **动手实践**: 不要只看文档，一定要运行示例代码
2. **修改代码**: 尝试修改参数，观察行为变化
3. **查看日志**: 日志包含大量有用信息
4. **使用 UI**: Web UI 可以直观地看到集群状态
5. **阅读源码**: Sarama 库的源码写得很好，值得学习
6. **循序渐进**: 从简单示例开始，逐步深入

## 📚 推荐资源

- [Apache Kafka 官方文档](https://kafka.apache.org/documentation/)
- [Sarama 文档](https://pkg.go.dev/github.com/IBM/sarama)
- [Confluent 博客](https://www.confluent.io/blog/)
- [Kafka 权威指南（书籍）](https://www.confluent.io/resources/kafka-the-definitive-guide/)

## 🤝 下一步

完成快速入门后，您可以：

1. 深入学习 [Kafka 架构](./docs/02-kafka-architecture.md)
2. 探索 [高级特性](./examples/)
3. 构建自己的实战项目
4. 学习 Kafka Streams 和 Kafka Connect

## ❓ 需要帮助？

- 查看 [README.md](./README.md) 了解项目结构
- 阅读各个示例的 README 文档
- 查看 [Kafka 官方文档](https://kafka.apache.org/)

祝学习愉快！🎉
