# 简单生产者示例

本示例演示如何使用 Sarama 创建一个简单的 Kafka 生产者，发送消息到指定的 Topic。

## 功能特性

- ✅ 同步发送消息
- ✅ 错误处理
- ✅ 消息确认
- ✅ 支持带 Key 的消息
- ✅ 优雅关闭

## 代码说明

### 核心配置

```go
config := sarama.NewConfig()
config.Producer.Return.Successes = true  // 等待成功响应
config.Producer.RequiredAcks = sarama.WaitForAll  // 等待所有 ISR 确认
config.Producer.Retry.Max = 3  // 失败重试 3 次
```

### 消息发送

```go
msg := &sarama.ProducerMessage{
    Topic: "example-topic",
    Key:   sarama.StringEncoder("user-123"),  // 可选：确保有序
    Value: sarama.StringEncoder("Hello Kafka"),
}

partition, offset, err := producer.SendMessage(msg)
```

## 运行示例

### 1. 启动 Kafka 集群

```bash
# 在项目根目录
docker-compose up -d

# 等待服务启动（约30秒）
docker-compose ps
```

### 2. 运行生产者

```bash
# 方式 1: 直接运行
go run examples/01-simple-producer/main.go

# 方式 2: 编译后运行
go build -o bin/simple-producer examples/01-simple-producer/main.go
./bin/simple-producer
```

### 3. 验证消息

```bash
# 方式 1: 使用 Kafka Console Consumer
docker exec -it kafka1 kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic example-topic \
  --from-beginning \
  --property print.key=true \
  --property print.timestamp=true

# 方式 2: 访问 Kafka UI
# 打开浏览器: http://localhost:8080
# 导航到 Topics -> example-topic -> Messages

# 方式 3: 使用 Kafdrop
# 打开浏览器: http://localhost:9000
```

## 输出示例

```
[2025-11-06 10:30:00] [INFO] [Producer] 连接到 Kafka 集群...
[2025-11-06 10:30:00] [INFO] [Producer] Broker 地址: [localhost:19092 localhost:29092 localhost:39092]
[2025-11-06 10:30:01] [INFO] [Producer] 成功连接到 Kafka
[2025-11-06 10:30:01] [INFO] [Producer] 开始发送消息...
[2025-11-06 10:30:01] [INFO] [Producer] 消息已发送 -> Topic: example-topic, Partition: 0, Offset: 0
[2025-11-06 10:30:01] [INFO] [Producer] 消息已发送 -> Topic: example-topic, Partition: 1, Offset: 0
[2025-11-06 10:30:01] [INFO] [Producer] 消息已发送 -> Topic: example-topic, Partition: 2, Offset: 0
[2025-11-06 10:30:02] [INFO] [Producer] 总共发送了 10 条消息
[2025-11-06 10:30:02] [INFO] [Producer] 关闭生产者
```

## 扩展练习

### 1. 发送 JSON 消息

```go
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

user := User{ID: "123", Name: "Alice", Email: "alice@example.com"}
jsonData, _ := json.Marshal(user)

msg := &sarama.ProducerMessage{
    Topic: "users",
    Value: sarama.ByteEncoder(jsonData),
}
```

### 2. 添加消息头

```go
msg := &sarama.ProducerMessage{
    Topic: "events",
    Headers: []sarama.RecordHeader{
        {Key: []byte("source"), Value: []byte("web-app")},
        {Key: []byte("version"), Value: []byte("1.0")},
    },
    Value: sarama.StringEncoder("event data"),
}
```

### 3. 指定分区

```go
msg := &sarama.ProducerMessage{
    Topic:     "orders",
    Partition: 2,  // 明确指定分区
    Value:     sarama.StringEncoder("order data"),
}
```

## 常见问题

### 1. 连接失败

```
错误: kafka: client has run out of available brokers
解决: 检查 Docker 容器是否正常运行
      docker-compose ps
```

### 2. 消息发送超时

```
错误: kafka: Failed to produce message to topic
解决: 增加超时时间
      config.Producer.Timeout = 10 * time.Second
```

### 3. Topic 不存在

```
错误: kafka: Unknown topic
解决: 确保启用了自动创建 Topic
      或手动创建: 
      docker exec kafka1 kafka-topics.sh --create \
        --bootstrap-server localhost:9092 \
        --topic example-topic \
        --partitions 3 \
        --replication-factor 2
```

## 下一步

- 💻 运行 [简单消费者示例](../02-simple-consumer/)
- 💻 学习 [消费者组](../03-consumer-group/)
- 💻 探索 [同步/异步生产者](../04-sync-async-producer/)
