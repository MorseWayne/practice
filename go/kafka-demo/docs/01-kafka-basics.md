# Kafka 核心概念

## 什么是 Apache Kafka？

Apache Kafka 是一个**分布式流处理平台**，最初由 LinkedIn 开发，后来贡献给了 Apache 软件基金会。它主要用于构建实时数据管道和流式应用。

### 核心特点

- **高吞吐量**: 单机支持数百万条消息/秒
- **可扩展性**: 可以水平扩展到数百个节点
- **持久性**: 消息持久化到磁盘，支持数据重放
- **容错性**: 支持数据副本，保证高可用
- **低延迟**: 毫秒级的消息传递延迟

## 核心概念

### 1. Topic (主题)

**Topic** 是 Kafka 中消息的分类单位，类似于数据库中的表或文件系统中的文件夹。

```
特点：
- 每个 Topic 可以有多个生产者和消费者
- 消息按照 Topic 进行组织
- Topic 是逻辑概念，物理上由多个 Partition 组成
```

**示例场景**:
- `user-events` - 用户行为事件
- `order-created` - 订单创建事件  
- `payment-notifications` - 支付通知
- `system-logs` - 系统日志

### 2. Partition (分区)

**Partition** 是 Topic 的物理分片，每个 Partition 是一个有序的、不可变的消息序列。

```
特点：
- 每个 Partition 内部消息有序
- 不同 Partition 之间消息无序
- 每个 Partition 可以有多个副本
- Partition 是 Kafka 并行处理的基本单位
```

**分区示意图**:
```
Topic: user-events
├── Partition 0: [msg1, msg4, msg7, ...]
├── Partition 1: [msg2, msg5, msg8, ...]
└── Partition 2: [msg3, msg6, msg9, ...]
```

**为什么需要分区？**
- **并行处理**: 多个消费者可以同时消费不同分区
- **负载均衡**: 消息分散到不同分区，提高吞吐量
- **可扩展性**: 可以通过增加分区数来提升性能

### 3. Producer (生产者)

**Producer** 负责向 Kafka 发送消息。

```
职责：
- 选择消息发送到哪个 Topic
- 决定消息发送到哪个 Partition
- 支持批量发送提高效率
```

**分区策略**:
1. **指定 Partition**: 直接指定消息发送到哪个分区
2. **基于 Key**: 相同 Key 的消息发送到同一分区
3. **轮询 (Round-robin)**: 均匀分配到各个分区

### 4. Consumer (消费者)

**Consumer** 负责从 Kafka 读取消息。

```
特点：
- 每个 Consumer 属于一个 Consumer Group
- 记录消费位置 (Offset)
- 支持从指定位置开始消费
- 可以重复消费历史消息
```

### 5. Consumer Group (消费者组)

**Consumer Group** 是一组协作消费的消费者集合。

```
核心规则：
- 同一个 Consumer Group 中，每个 Partition 只能被一个 Consumer 消费
- 不同 Consumer Group 可以独立消费同一个 Topic
- 通过增加 Consumer 实现负载均衡
```

**负载均衡示例**:
```
Topic: orders (3 个分区)
Consumer Group A: (2 个消费者)
├── Consumer-1: 消费 Partition 0, 1
└── Consumer-2: 消费 Partition 2

Consumer Group B: (3 个消费者)  
├── Consumer-3: 消费 Partition 0
├── Consumer-4: 消费 Partition 1
└── Consumer-5: 消费 Partition 2
```

### 6. Broker (代理)

**Broker** 是 Kafka 集群中的服务器节点。

```
职责：
- 存储消息数据
- 处理生产者和消费者请求
- 管理 Partition 副本
- 参与 Leader 选举
```

### 7. Offset (偏移量)

**Offset** 是消息在 Partition 中的唯一标识，是一个递增的整数。

```
类型：
- Current Offset: 当前消费到的位置
- Committed Offset: 已提交的消费位置
- Log End Offset: 分区中最新消息的位置
```

**Offset 管理**:
```
Partition 0: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
             ^        ^        ^
          消费起点   已提交   最新消息
```

### 8. Replication (副本)

**副本机制**保证数据的可靠性和高可用性。

```
角色：
- Leader: 处理所有读写请求
- Follower: 从 Leader 同步数据，作为备份

副本分配：
- 每个 Partition 可以配置多个副本
- 副本分布在不同的 Broker 上
- 当 Leader 失效时，从 Follower 中选举新 Leader
```

## Kafka 消息结构

### 消息组成

```go
type Message struct {
    Key       []byte    // 消息键 (可选，用于分区)
    Value     []byte    // 消息内容
    Timestamp time.Time // 时间戳
    Headers   []Header  // 消息头 (元数据)
    Partition int32     // 分区号
    Offset    int64     // 偏移量
}
```

### 消息示例

```json
{
  "key": "user-123",
  "value": {
    "event": "login",
    "userId": "123",
    "timestamp": "2025-11-06T10:30:00Z",
    "ip": "192.168.1.1"
  },
  "headers": {
    "source": "web-app",
    "version": "1.0"
  }
}
```

## Kafka 使用场景

### 1. 消息队列
替代传统消息队列（RabbitMQ、ActiveMQ），提供更高的吞吐量。

### 2. 日志聚合
收集各个服务的日志，集中存储和分析。

### 3. 流处理
实时处理数据流，如实时统计、监控告警等。

### 4. 事件溯源
记录系统中的所有事件，支持事件回放和审计。

### 5. 数据管道
作为数据管道的中间层，连接不同的数据系统。

### 6. 微服务通信
服务之间的异步通信和事件驱动架构。

## 消息传递语义

### 1. At Most Once (最多一次)
- 消息可能丢失，但不会重复
- 适用于对数据丢失不敏感的场景

### 2. At Least Once (至少一次)
- 消息不会丢失，但可能重复
- Kafka 默认保证的语义
- 需要消费者实现幂等性

### 3. Exactly Once (精确一次)
- 消息既不丢失也不重复
- 通过事务和幂等生产者实现
- 性能开销较大

## 关键配置参数

### Producer 关键参数

```properties
# 确认机制
acks=all              # 0/1/all，控制消息可靠性

# 重试配置
retries=3             # 发送失败重试次数
retry.backoff.ms=100  # 重试间隔

# 批处理
batch.size=16384      # 批次大小（字节）
linger.ms=10          # 等待时间（毫秒）

# 压缩
compression.type=gzip # none/gzip/snappy/lz4/zstd
```

### Consumer 关键参数

```properties
# 消费者组
group.id=my-group     # 消费者组ID

# 自动提交
enable.auto.commit=true     # 是否自动提交offset
auto.commit.interval.ms=5000 # 自动提交间隔

# 消费位置
auto.offset.reset=latest    # earliest/latest/none

# 拉取配置
fetch.min.bytes=1           # 最小拉取字节数
fetch.max.wait.ms=500       # 最大等待时间
```

## 下一步

- 📖 阅读 [Kafka 架构设计](./02-kafka-architecture.md)
- 📖 阅读 [消息模型与分区](./03-message-partition.md)
- 🚀 开始 [环境搭建](./04-setup-environment.md)
