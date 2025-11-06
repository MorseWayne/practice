package config

import (
	"time"

	"github.com/IBM/sarama"
)

// KafkaConfig Kafka 配置
type KafkaConfig struct {
	Brokers []string
	Version sarama.KafkaVersion
}

// DefaultKafkaConfig 默认 Kafka 配置
func DefaultKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Brokers: []string{
			"localhost:19092",
			"localhost:29092",
			"localhost:39092",
		},
		Version: sarama.V3_6_0_0,
	}
}

// NewProducerConfig 创建生产者配置
func NewProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0

	// 生产者配置
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认
	config.Producer.Retry.Max = 3
	config.Producer.Compression = sarama.CompressionSnappy

	// 幂等性配置（防止重复）
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	return config
}

// NewConsumerConfig 创建消费者配置
func NewConsumerConfig(groupID string) *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0

	// 消费者配置
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky

	// 手动提交 Offset
	config.Consumer.Offsets.AutoCommit.Enable = false

	// 会话超时配置
	config.Consumer.Group.Session.Timeout = 20 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 6 * time.Second

	// 拉取配置
	config.Consumer.Fetch.Min = 1024            // 1KB
	config.Consumer.Fetch.Default = 1024 * 1024 // 1MB

	return config
}

// NewAsyncProducerConfig 创建异步生产者配置
func NewAsyncProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0

	// 异步生产者配置
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForLocal // 只等待 Leader 确认
	config.Producer.Compression = sarama.CompressionSnappy

	// 批处理配置
	config.Producer.Flush.Messages = 100
	config.Producer.Flush.Frequency = 10 * time.Millisecond
	config.Producer.Flush.MaxMessages = 1000

	// 异步模式下可以开启管道
	config.Net.MaxOpenRequests = 5

	return config
}
