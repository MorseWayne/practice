# 简单消费者示例

本示例演示如何使用 Sarama 创建一个简单的 Kafka 消费者，从指定的 Topic 消费消息。

## 功能特性

- ✅ 消费者组订阅
- ✅ 自动 Rebalance
- ✅ 手动提交 Offset
- ✅ 错误处理
- ✅ 优雅关闭

## 代码说明

### 消费者组配置

```go
config := sarama.NewConfig()
config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
config.Consumer.Offsets.AutoCommit.Enable = false  // 手动提交
config.Consumer.Offsets.Initial = sarama.OffsetNewest  // 从最新位置开始
```

### 消息处理

```go
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for message := range claim.Messages() {
        // 处理消息
        processMessage(message)
        
        // 标记消息为已处理
        session.MarkMessage(message, "")
    }
    return nil
}
```

## 运行示例

### 1. 确保 Kafka 集群运行

```bash
docker-compose ps
```

### 2. 先运行生产者（发送消息）

```bash
go run examples/01-simple-producer/main.go
```

### 3. 运行消费者

```bash
# 方式 1: 直接运行
go run examples/02-simple-consumer/main.go

# 方式 2: 编译后运行
go build -o bin/simple-consumer examples/02-simple-consumer/main.go
./bin/simple-consumer
```

### 4. 测试多消费者负载均衡

在不同终端窗口运行多个消费者实例：

```bash
# 终端 1
go run examples/02-simple-consumer/main.go

# 终端 2
go run examples/02-simple-consumer/main.go

# 终端 3
go run examples/02-simple-consumer/main.go
```

每个消费者会被分配不同的分区，实现负载均衡。

## 输出示例

```
[2025-11-06 10:35:00] [Consumer] 启动消费者...
[2025-11-06 10:35:00] [Consumer] 消费者组: example-consumer-group
[2025-11-06 10:35:00] [Consumer] 订阅 Topic: example-topic
[2025-11-06 10:35:01] [Consumer] 成功加入消费者组
[2025-11-06 10:35:01] [Consumer] 分配的分区: [0 1 2]
[2025-11-06 10:35:02] [Consumer] 📨 收到消息:
  Topic: example-topic
  Partition: 0
  Offset: 5
  Key: key-1
  Value: 消息 #1 - 时间: 2025-11-06T10:30:01Z
  Timestamp: 2025-11-06 10:30:01
[2025-11-06 10:35:02] [Consumer] ✅ 消息处理完成，已提交 Offset
```

## 消费者组行为

### 单个消费者

```
Topic: example-topic (3 个分区)
Consumer Group: my-group

Consumer-1 消费:
├── Partition 0
├── Partition 1
└── Partition 2
```

### 两个消费者

```
Topic: example-topic (3 个分区)
Consumer Group: my-group

Consumer-1 消费:
├── Partition 0
└── Partition 1

Consumer-2 消费:
└── Partition 2
```

### 三个消费者

```
Topic: example-topic (3 个分区)
Consumer Group: my-group

Consumer-1 消费: Partition 0
Consumer-2 消费: Partition 1
Consumer-3 消费: Partition 2
```

### 四个消费者（消费者过多）

```
Topic: example-topic (3 个分区)
Consumer Group: my-group

Consumer-1 消费: Partition 0
Consumer-2 消费: Partition 1
Consumer-3 消费: Partition 2
Consumer-4 消费: （空闲，无分区分配）
```

## Offset 管理

### 查看消费者组状态

```bash
docker exec -it kafka1 kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --group example-consumer-group \
  --describe
```

输出示例：

```
GROUP                    TOPIC           PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG
example-consumer-group   example-topic   0          10              10              0
example-consumer-group   example-topic   1          11              11              0
example-consumer-group   example-topic   2          9               9               0
```

### 重置 Offset

```bash
# 重置到最早
docker exec -it kafka1 kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --group example-consumer-group \
  --topic example-topic \
  --reset-offsets \
  --to-earliest \
  --execute

# 重置到最新
docker exec -it kafka1 kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --group example-consumer-group \
  --topic example-topic \
  --reset-offsets \
  --to-latest \
  --execute

# 重置到指定位置
docker exec -it kafka1 kafka-consumer-groups.sh \
  --bootstrap-server localhost:9092 \
  --group example-consumer-group \
  --topic example-topic \
  --reset-offsets \
  --to-offset 5 \
  --execute
```

## 扩展练习

### 1. 处理 JSON 消息

```go
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func processMessage(message *sarama.ConsumerMessage) {
    var user User
    if err := json.Unmarshal(message.Value, &user); err != nil {
        log.Printf("JSON 解析失败: %v", err)
        return
    }
    log.Printf("处理用户: %+v", user)
}
```

### 2. 批量提交 Offset

```go
batch := make([]*sarama.ConsumerMessage, 0, 100)
for message := range claim.Messages() {
    batch = append(batch, message)
    
    if len(batch) >= 100 {
        // 批量处理
        processBatch(batch)
        
        // 提交最后一条消息的 Offset
        session.MarkMessage(batch[len(batch)-1], "")
        session.Commit()
        
        batch = batch[:0]
    }
}
```

### 3. 错误重试

```go
func processMessage(message *sarama.ConsumerMessage) error {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        err := doProcess(message)
        if err == nil {
            return nil
        }
        log.Printf("处理失败，重试 %d/%d: %v", i+1, maxRetries, err)
        time.Sleep(time.Second * time.Duration(i+1))
    }
    return fmt.Errorf("处理失败，已达最大重试次数")
}
```

## 常见问题

### 1. Rebalance 频繁发生

```
原因: 消息处理时间过长，超过会话超时
解决: 增加会话超时时间
      config.Consumer.Group.Session.Timeout = 30 * time.Second
      config.Consumer.MaxProcessingTime = 5 * time.Minute
```

### 2. 消息重复消费

```
原因: 消息处理完但 Offset 未提交前程序崩溃
解决: 1. 实现幂等性处理
      2. 使用事务
      3. 及时提交 Offset
```

### 3. 消息丢失

```
原因: 自动提交 Offset，但消息处理失败
解决: 使用手动提交，确保消息处理成功后再提交
      config.Consumer.Offsets.AutoCommit.Enable = false
```

### 4. 消费延迟（Lag）过大

```
原因: 消费速度慢于生产速度
解决: 1. 增加消费者数量
      2. 增加分区数
      3. 优化消息处理逻辑
      4. 使用批量处理
```

## 下一步

- 💻 学习 [消费者组详解](../03-consumer-group/)
- 💻 探索 [同步/异步生产者](../04-sync-async-producer/)
- 💻 学习 [批量处理](../06-batch-processing/)
