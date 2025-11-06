# 消息模型与分区策略

## Kafka 消息模型

### 消息结构详解

#### 完整消息格式

```go
type ProducerRecord struct {
    Topic      string              // 目标主题
    Partition  *int32              // 目标分区（可选）
    Key        []byte              // 消息键
    Value      []byte              // 消息值
    Headers    []RecordHeader      // 消息头
    Timestamp  time.Time           // 时间戳
}

type ConsumerRecord struct {
    Topic      string              // 主题
    Partition  int32               // 分区
    Offset     int64               // 偏移量
    Key        []byte              // 键
    Value      []byte              // 值
    Headers    []RecordHeader      // 消息头
    Timestamp  time.Time           // 时间戳
}

type RecordHeader struct {
    Key   string
    Value []byte
}
```

#### 消息示例

```json
{
  "topic": "user-events",
  "partition": 2,
  "offset": 12345,
  "timestamp": "2025-11-06T10:30:00.000Z",
  "key": "user-123",
  "value": {
    "eventType": "login",
    "userId": "123",
    "deviceId": "mobile-001",
    "location": "Beijing",
    "ip": "192.168.1.100"
  },
  "headers": [
    {"key": "source", "value": "mobile-app"},
    {"key": "version", "value": "2.0"},
    {"key": "traceId", "value": "abc-123-xyz"}
  ]
}
```

### Key 的作用

#### 1. 分区路由

```
相同 Key 的消息会被发送到同一个 Partition

示例：
Key = "user-123" -> hash(key) % partitions = Partition 2
Key = "user-456" -> hash(key) % partitions = Partition 1
Key = "user-123" -> hash(key) % partitions = Partition 2  ✓ 相同

好处：
- 保证相同实体的消息有序
- 便于状态管理
- 支持本地聚合
```

#### 2. 日志压缩 (Log Compaction)

```
启用日志压缩时，Kafka 只保留每个 Key 的最新值

Before Compaction:
offset: 0  key: A  value: 1
offset: 1  key: B  value: 2
offset: 2  key: A  value: 3  <- 新值
offset: 3  key: C  value: 4
offset: 4  key: B  value: 5  <- 新值

After Compaction:
offset: 2  key: A  value: 3  <- 保留
offset: 3  key: C  value: 4
offset: 4  key: B  value: 5  <- 保留

应用场景：
- 数据库变更日志 (CDC)
- 配置管理
- 用户状态快照
```

### Headers 的用途

```go
// 1. 链路追踪
headers := []sarama.RecordHeader{
    {Key: []byte("trace-id"), Value: []byte("abc-123")},
    {Key: []byte("span-id"), Value: []byte("xyz-789")},
}

// 2. 消息来源
headers = append(headers, sarama.RecordHeader{
    Key:   []byte("source"),
    Value: []byte("order-service"),
})

// 3. 消息类型
headers = append(headers, sarama.RecordHeader{
    Key:   []byte("event-type"),
    Value: []byte("OrderCreated"),
})

// 4. 版本控制
headers = append(headers, sarama.RecordHeader{
    Key:   []byte("schema-version"),
    Value: []byte("v2.0"),
})
```

## 分区策略

### 1. 默认分区器 (Default Partitioner)

```
规则：
├── 如果指定了 Partition -> 使用指定的分区
├── 如果提供了 Key -> hash(key) % numPartitions
└── 如果没有 Key -> 轮询 (Round-robin) 或 Sticky

代码示例：
// 指定分区
msg := &sarama.ProducerMessage{
    Topic:     "orders",
    Partition: 2,  // 明确指定分区 2
    Value:     sarama.StringEncoder("order data"),
}

// 基于 Key
msg := &sarama.ProducerMessage{
    Topic: "orders",
    Key:   sarama.StringEncoder("user-123"),  // 相同 key 总是到同一分区
    Value: sarama.StringEncoder("order data"),
}

// 轮询
msg := &sarama.ProducerMessage{
    Topic: "orders",
    // 没有指定 Key 和 Partition，使用轮询
    Value: sarama.StringEncoder("order data"),
}
```

### 2. 自定义分区器

#### 示例 1: 按地区分区

```go
type RegionPartitioner struct{}

func (p *RegionPartitioner) Partition(
    message *sarama.ProducerMessage,
    numPartitions int32,
) (int32, error) {
    // 从消息头中获取地区信息
    region := ""
    for _, header := range message.Headers {
        if string(header.Key) == "region" {
            region = string(header.Value)
            break
        }
    }
    
    // 根据地区分配分区
    switch region {
    case "north":
        return 0, nil
    case "south":
        return 1, nil
    case "east":
        return 2, nil
    case "west":
        return 3, nil
    default:
        // 默认使用轮询
        return rand.Int31n(numPartitions), nil
    }
}

func (p *RegionPartitioner) RequiresConsistency() bool {
    return false
}
```

#### 示例 2: 按优先级分区

```go
type PriorityPartitioner struct{}

func (p *PriorityPartitioner) Partition(
    message *sarama.ProducerMessage,
    numPartitions int32,
) (int32, error) {
    // 假设消息头包含优先级
    priority := "normal"
    for _, header := range message.Headers {
        if string(header.Key) == "priority" {
            priority = string(header.Value)
            break
        }
    }
    
    // 高优先级消息到专用分区
    if priority == "high" {
        return 0, nil  // 分区 0 专门处理高优先级
    }
    
    // 其他消息均匀分布到其他分区
    partition := 1 + (rand.Int31n(numPartitions - 1))
    return partition, nil
}

func (p *PriorityPartitioner) RequiresConsistency() bool {
    return false
}
```

#### 示例 3: 基于时间的分区

```go
type TimeBasedPartitioner struct{}

func (p *TimeBasedPartitioner) Partition(
    message *sarama.ProducerMessage,
    numPartitions int32,
) (int32, error) {
    // 根据小时数分区，适合时序数据
    hour := time.Now().Hour()
    partition := int32(hour % int(numPartitions))
    return partition, nil
}

func (p *TimeBasedPartitioner) RequiresConsistency() bool {
    return false
}
```

### 3. 分区数选择

#### 考虑因素

```
1. 吞吐量需求
   - 更多分区 = 更高并行度 = 更高吞吐量
   - 每个分区可以被不同的消费者处理

2. 消费者数量
   - 消费者数量 ≤ 分区数
   - 过多分区会浪费（消费者闲置）

3. 数据有序性要求
   - 需要全局有序 -> 1 个分区
   - 需要局部有序 -> 基于 Key 分区

4. 延迟要求
   - 更多分区 = 更多文件 = 更高延迟
   - Leader 选举时间也会增加

5. 存储容量
   - 每个分区占用磁盘空间
   - 需要考虑副本因子
```

#### 经验法则

```
推荐公式：
分区数 = max(
    目标吞吐量 / 单分区吞吐量,
    消费者数量
)

示例：
- 目标吞吐量: 100 MB/s
- 单分区吞吐量: 10 MB/s
- 消费者数量: 5

分区数 = max(100/10, 5) = max(10, 5) = 10

建议范围：
- 小规模: 3-10 个分区
- 中等规模: 10-50 个分区
- 大规模: 50-200 个分区
- 极限: 不超过 4000 个分区/broker
```

## Offset 管理

### Offset 类型

```
1. Current Offset (当前偏移量)
   - 消费者当前读取到的位置

2. Committed Offset (已提交偏移量)
   - 消费者已确认处理完成的位置
   - 存储在 __consumer_offsets topic

3. Log End Offset (LEO)
   - 分区中下一条消息将写入的位置
   - 即当前最大 offset + 1

4. High Water Mark (HWM)
   - ISR 中所有副本都已同步到的位置
   - 消费者只能读取到 HWM 之前的消息
```

### Offset 提交策略

#### 1. 自动提交

```go
config := sarama.NewConfig()
config.Consumer.Offsets.AutoCommit.Enable = true
config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

优点：
- 使用简单
- 无需手动管理

缺点：
- 可能导致消息丢失（先提交后处理）
- 可能导致重复消费（处理慢但已提交）
```

#### 2. 手动同步提交

```go
config := sarama.NewConfig()
config.Consumer.Offsets.AutoCommit.Enable = false

// 处理消息后提交
for message := range consumer.Messages() {
    // 处理消息
    processMessage(message)
    
    // 同步提交
    session.MarkMessage(message, "")
    session.Commit()
}

优点：
- 精确控制提交时机
- 保证消息被处理后才提交

缺点：
- 性能较低（同步等待）
```

#### 3. 手动异步提交

```go
// 批量处理后异步提交
batch := make([]*sarama.ConsumerMessage, 0, 100)
for message := range consumer.Messages() {
    batch = append(batch, message)
    
    if len(batch) >= 100 {
        // 批量处理
        processBatch(batch)
        
        // 标记最后一条消息
        session.MarkMessage(batch[len(batch)-1], "")
        
        batch = batch[:0]
    }
}

优点：
- 平衡性能和可靠性
- 减少提交次数

缺点：
- 复杂度较高
- 需要处理失败场景
```

### Offset 重置

```go
// 从最早位置开始消费
config.Consumer.Offsets.Initial = sarama.OffsetOldest

// 从最新位置开始消费
config.Consumer.Offsets.Initial = sarama.OffsetNewest

// 手动设置 Offset
consumer.Seek(partition, offset)

// 根据时间戳查找 Offset
timestamp := time.Now().Add(-24 * time.Hour).Unix()
// 使用 Kafka Admin API 查找
```

## 消息顺序保证

### 场景 1: 全局有序

```
要求：所有消息严格有序

方案：
- 只使用 1 个分区
- max.in.flight.requests.per.connection = 1

代码：
config := sarama.NewConfig()
config.Producer.Idempotent = true
config.Net.MaxOpenRequests = 1

缺点：
- 吞吐量受限
- 无法并行消费
```

### 场景 2: 局部有序（推荐）

```
要求：同一实体的消息有序

方案：
- 使用相同的 Key
- 消息路由到同一分区

代码：
// 订单消息使用订单 ID 作为 Key
msg := &sarama.ProducerMessage{
    Topic: "orders",
    Key:   sarama.StringEncoder(order.ID),  // 相同订单 ID -> 同一分区
    Value: orderJSON,
}

优点：
- 保证关键业务有序
- 可以并行处理不同实体
- 吞吐量高
```

### 场景 3: 无序（最高性能）

```
要求：不要求顺序

方案：
- 不指定 Key
- 使用异步发送
- 允许重试

代码：
config := sarama.NewConfig()
config.Producer.RequiredAcks = sarama.WaitForLocal
config.Producer.Retry.Max = 3
config.Producer.Return.Successes = true

优点：
- 最高吞吐量
- 最低延迟
```

## 消息过期与清理

### 1. 基于时间的保留

```properties
# 保留 7 天
log.retention.hours=168

# 或使用分钟
log.retention.minutes=10080

# 或使用毫秒
log.retention.ms=604800000
```

### 2. 基于大小的保留

```properties
# 每个分区最大 1GB
log.retention.bytes=1073741824

# 注意：两个条件满足任意一个就会清理
```

### 3. 日志压缩 (Log Compaction)

```properties
# 启用日志压缩
log.cleanup.policy=compact

# 或同时使用删除和压缩
log.cleanup.policy=compact,delete

配置：
log.cleaner.enable=true
log.cleaner.min.compaction.lag.ms=0
log.cleaner.max.compaction.lag.ms=86400000

应用场景：
- 数据库变更日志
- 用户配置快照
- 缓存更新
```

### 日志压缩工作原理

```
原始日志：
offset  key   value
0       A     v1
1       B     v1
2       A     v2    <- A 的新值
3       C     v1
4       B     v2    <- B 的新值
5       A     v3    <- A 的最新值
6       D     v1

压缩后：
offset  key   value
2       A     v2    <- 删除 offset 0
3       C     v1
4       B     v2    <- 删除 offset 1
5       A     v3    <- 保留最新，删除 offset 2
6       D     v1

最终结果：每个 key 只保留最新的值
```

## 分区再平衡优化

### 减少 Rebalance 影响

```go
config := sarama.NewConfig()

// 1. 增加会话超时时间
config.Consumer.Group.Session.Timeout = 20 * time.Second

// 2. 增加心跳间隔
config.Consumer.Group.Heartbeat.Interval = 6 * time.Second

// 3. 增加处理时间限制
config.Consumer.MaxProcessingTime = 10 * time.Minute

// 4. 使用 Sticky 分配策略
config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
```

### 优雅关闭

```go
// 捕获信号
signals := make(chan os.Signal, 1)
signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

go func() {
    <-signals
    // 优雅关闭消费者
    if err := consumer.Close(); err != nil {
        log.Printf("Error closing consumer: %v", err)
    }
}()
```

## 性能调优建议

### Producer 调优

```go
config := sarama.NewConfig()

// 1. 批处理
config.Producer.Flush.Messages = 100
config.Producer.Flush.Frequency = 10 * time.Millisecond

// 2. 压缩
config.Producer.Compression = sarama.CompressionSnappy

// 3. 异步发送
config.Producer.Return.Successes = true
config.Producer.Return.Errors = true

// 4. 缓冲区大小
config.ChannelBufferSize = 256
```

### Consumer 调优

```go
config := sarama.NewConfig()

// 1. 批量拉取
config.Consumer.Fetch.Min = 1024      // 1 KB
config.Consumer.Fetch.Default = 1048576  // 1 MB
config.Consumer.Fetch.Max = 52428800  // 50 MB

// 2. 等待时间
config.Consumer.MaxWaitTime = 500 * time.Millisecond

// 3. 并发处理
// 在消费逻辑中使用 goroutine 池
```

## 下一步

- 🚀 开始 [环境搭建](./04-setup-environment.md)
- 💻 运行 [简单生产者示例](../examples/01-simple-producer/)
- 💻 运行 [简单消费者示例](../examples/02-simple-consumer/)
