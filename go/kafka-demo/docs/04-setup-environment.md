# 环境搭建指南

本指南将帮助您使用 Docker Compose 快速搭建 Kafka 开发环境。

## 前置条件

- Docker 20.10+
- Docker Compose 1.29+
- 至少 4GB 可用内存

## Docker Compose 配置

我们将创建一个包含以下组件的 Kafka 集群：

- **ZooKeeper**: 1 个节点（协调服务）
- **Kafka Broker**: 3 个节点（高可用集群）
- **Kafka UI**: Web 管理界面
- **Kafdrop**: 消息查看工具

## 快速启动

### 1. 启动所有服务

```bash
# 启动集群（后台运行）
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f kafka1
```

### 2. 验证集群

```bash
# 进入 Kafka 容器
docker exec -it kafka1 bash

# 创建测试 Topic
kafka-topics.sh --create \
  --bootstrap-server localhost:9092 \
  --topic test-topic \
  --partitions 3 \
  --replication-factor 2

# 查看 Topic 列表
kafka-topics.sh --list \
  --bootstrap-server localhost:9092

# 查看 Topic 详情
kafka-topics.sh --describe \
  --bootstrap-server localhost:9092 \
  --topic test-topic

# 发送测试消息
echo "Hello Kafka" | kafka-console-producer.sh \
  --bootstrap-server localhost:9092 \
  --topic test-topic

# 消费测试消息
kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic test-topic \
  --from-beginning
```

### 3. 访问 Web UI

- **Kafka UI**: http://localhost:8080
- **Kafdrop**: http://localhost:9000

这些工具提供了可视化的界面来：
- 查看集群状态
- 管理 Topic
- 浏览消息
- 监控消费者组

### 4. 停止服务

```bash
# 停止所有服务
docker-compose stop

# 停止并删除容器（保留数据卷）
docker-compose down

# 停止并删除所有内容（包括数据）
docker-compose down -v
```

## 连接配置

### Go 客户端连接

```go
// 开发环境配置
config := sarama.NewConfig()
config.Version = sarama.V3_6_0_0

brokers := []string{
    "localhost:19092",  // kafka1
    "localhost:29092",  // kafka2
    "localhost:39092",  // kafka3
}

// 创建生产者
producer, err := sarama.NewSyncProducer(brokers, config)
if err != nil {
    log.Fatal(err)
}
defer producer.Close()
```

### 端口说明

```
服务端口映射：
├── ZooKeeper
│   └── 2181:2181       # 客户端连接
├── Kafka Broker 1
│   ├── 19092:9092      # 外部访问（从 Docker 外）
│   └── 9092            # 内部访问（容器间）
├── Kafka Broker 2
│   ├── 29092:9092      # 外部访问
│   └── 9092            # 内部访问
├── Kafka Broker 3
│   ├── 39092:9092      # 外部访问
│   └── 9092            # 内部访问
├── Kafka UI
│   └── 8080:8080       # Web 界面
└── Kafdrop
    └── 9000:9000       # Web 界面
```

## 常用命令

### Topic 管理

```bash
# 进入 Kafka 容器
docker exec -it kafka1 bash

# 创建 Topic
kafka-topics.sh --create \
  --bootstrap-server localhost:9092 \
  --topic my-topic \
  --partitions 3 \
  --replication-factor 2

# 列出所有 Topic
kafka-topics.sh --list \
  --bootstrap-server localhost:9092

# 查看 Topic 详情
kafka-topics.sh --describe \
  --bootstrap-server localhost:9092 \
  --topic my-topic

# 修改 Topic 分区数
kafka-topics.sh --alter \
  --bootstrap-server localhost:9092 \
  --topic my-topic \
  --partitions 5

# 删除 Topic
kafka-topics.sh --delete \
  --bootstrap-server localhost:9092 \
  --topic my-topic
```

### 消息生产与消费

```bash
# 生产消息（交互式）
kafka-console-producer.sh \
  --bootstrap-server localhost:9092 \
  --topic my-topic

# 带 Key 的生产者
kafka-console-producer.sh \
  --bootstrap-server localhost:9092 \
  --topic my-topic \
  --property "parse.key=true" \
  --property "key.separator=:"

# 消费消息（从最新）
kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic my-topic

# 消费消息（从头开始）
kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic my-topic \
  --from-beginning

# 显示 Key 和时间戳
kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic my-topic \
  --from-beginning \
  --property print.key=true \
  --property print.timestamp=true
```

### 消费者组管理

```bash
# 列出所有消费者组
kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --list

# 查看消费者组详情
kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --group my-group \
  --describe

# 重置消费者组 Offset（到最早）
kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --group my-group \
  --topic my-topic \
  --reset-offsets \
  --to-earliest \
  --execute

# 重置到指定位置
kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --group my-group \
  --topic my-topic \
  --reset-offsets \
  --to-offset 100 \
  --execute

# 删除消费者组（需先停止所有消费者）
kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --group my-group \
  --delete
```

### 性能测试

```bash
# 生产者性能测试
kafka-producer-perf-test.sh \
  --topic perf-test \
  --num-records 1000000 \
  --record-size 1024 \
  --throughput -1 \
  --producer-props \
    bootstrap.servers=localhost:9092 \
    acks=all

# 消费者性能测试
kafka-consumer-perf-test.sh \
  --bootstrap-server localhost:9092 \
  --topic perf-test \
  --messages 1000000 \
  --threads 1
```

## 故障排查

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs

# 跟踪实时日志
docker-compose logs -f kafka1

# 查看最后 100 行
docker-compose logs --tail=100 kafka1
```

### 检查服务健康状态

```bash
# 检查容器状态
docker-compose ps

# 检查 ZooKeeper 连接
echo stat | nc localhost 2181

# 检查 Kafka Broker
docker exec kafka1 kafka-broker-api-versions.sh \
  --bootstrap-server localhost:9092
```

### 常见问题

#### 1. Kafka 无法连接

```bash
# 检查网络
docker network inspect kafka-demo_default

# 检查端口监听
docker exec kafka1 netstat -tulpn | grep 9092

# 检查防火墙规则
sudo ufw status
```

#### 2. ZooKeeper 连接失败

```bash
# 验证 ZooKeeper 运行
docker exec zookeeper zkServer.sh status

# 测试连接
docker exec zookeeper zkCli.sh ls /
```

#### 3. 内存不足

```yaml
# 在 docker-compose.yml 中调整内存限制
environment:
  KAFKA_HEAP_OPTS: "-Xmx512M -Xms512M"
```

#### 4. 磁盘空间不足

```bash
# 清理旧日志
docker exec kafka1 kafka-configs.sh \
  --bootstrap-server localhost:9092 \
  --entity-type topics \
  --entity-name my-topic \
  --alter \
  --add-config retention.ms=3600000  # 1 小时

# 清理 Docker 资源
docker system prune -a --volumes
```

## 生产环境注意事项

### 1. 资源配置

```yaml
# 建议的生产环境配置
deploy:
  resources:
    limits:
      cpus: '2'
      memory: 4G
    reservations:
      cpus: '1'
      memory: 2G
```

### 2. 数据持久化

```yaml
# 使用命名卷
volumes:
  - kafka-data:/var/lib/kafka/data

# 定义顶层卷
volumes:
  kafka-data:
    driver: local
```

### 3. 网络配置

```yaml
# 使用自定义网络
networks:
  kafka-net:
    driver: bridge
    ipam:
      config:
        - subnet: 172.25.0.0/16
```

### 4. 安全配置

```yaml
# 启用 SASL 认证
environment:
  KAFKA_SECURITY_PROTOCOL: SASL_PLAINTEXT
  KAFKA_SASL_MECHANISM: PLAIN
  KAFKA_SASL_ENABLED_MECHANISMS: PLAIN
```

### 5. 监控配置

```yaml
# 启用 JMX 监控
environment:
  KAFKA_JMX_PORT: 9999
  KAFKA_JMX_HOSTNAME: kafka1
ports:
  - "9999:9999"
```

## 下一步

- 💻 运行 [简单生产者示例](../examples/01-simple-producer/)
- 💻 运行 [简单消费者示例](../examples/02-simple-consumer/)
- 💻 学习 [消费者组示例](../examples/03-consumer-group/)
