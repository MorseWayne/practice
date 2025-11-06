# Kafka + Go 学习指南

这是一个完整的 Kafka 学习项目，结合 Go 语言的实战示例，从基础到高级逐步深入。

## 📚 学习目录

### 1. Kafka 基础概念
- [Kafka 核心概念](./docs/01-kafka-basics.md)
- [Kafka 架构设计](./docs/02-kafka-architecture.md)
- [消息模型与分区](./docs/03-message-partition.md)

### 2. 环境搭建
- [使用 Docker Compose 快速搭建 Kafka 集群](./docs/04-setup-environment.md)

### 3. Go 基础示例
- [简单生产者 (Simple Producer)](./examples/01-simple-producer/)
- [简单消费者 (Simple Consumer)](./examples/02-simple-consumer/)
- [消费者组 (Consumer Group)](./examples/03-consumer-group/)

### 4. Go 高级特性
- [同步/异步生产者](./examples/04-sync-async-producer/)
- [消息事务 (Transactions)](./examples/05-transactions/)
- [批量发送与消费](./examples/06-batch-processing/)
- [消息拦截器与序列化](./examples/07-interceptors-serialization/)

### 5. 实战案例
- [订单处理系统](./examples/08-order-processing/)
- [日志聚合系统](./examples/09-log-aggregation/)
- [实时数据管道](./examples/10-data-pipeline/)

## 🚀 快速开始

### 前置条件
- Go 1.21+ 
- Docker & Docker Compose
- 基本的命令行操作知识

### 启动 Kafka 环境

```bash
# 启动 Kafka 集群
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f kafka
```

### 运行示例

```bash
# 初始化项目
go mod download

# 运行简单生产者
go run examples/01-simple-producer/main.go

# 运行简单消费者
go run examples/02-simple-consumer/main.go
```

## 📖 使用的 Go Kafka 客户端

本项目使用 [Sarama](https://github.com/IBM/sarama) - 一个纯 Go 实现的 Apache Kafka 客户端库。

也包含了使用 [kafka-go](https://github.com/segmentio/kafka-go) 的示例对比。

## 📝 学习路径建议

1. **第一天**: 阅读 Kafka 基础概念文档，理解核心术语
2. **第二天**: 搭建环境，运行简单的生产者和消费者示例
3. **第三天**: 学习消费者组，理解分区与负载均衡
4. **第四天**: 探索高级特性：事务、批处理等
5. **第五天**: 研究实战案例，理解实际应用场景

## 🔍 项目结构

```
kafka-demo/
├── docs/                    # 学习文档
├── examples/                # 代码示例
│   ├── 01-simple-producer/
│   ├── 02-simple-consumer/
│   └── ...
├── pkg/                     # 共享工具包
│   ├── config/             # 配置管理
│   └── logger/             # 日志工具
├── docker-compose.yml       # Docker 环境配置
├── go.mod
├── go.sum
└── README.md
```

## 🤝 最佳实践

- ✅ 总是正确处理错误
- ✅ 使用消费者组实现负载均衡
- ✅ 合理设置批处理大小
- ✅ 根据场景选择同步/异步发送
- ✅ 实现优雅关闭
- ✅ 监控生产和消费延迟

## 📚 参考资源

- [Apache Kafka 官方文档](https://kafka.apache.org/documentation/)
- [Sarama 文档](https://pkg.go.dev/github.com/IBM/sarama)
- [Kafka: The Definitive Guide](https://www.confluent.io/resources/kafka-the-definitive-guide/)

## 📄 License

MIT License
