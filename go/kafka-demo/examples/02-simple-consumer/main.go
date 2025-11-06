package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

const (
	topic         = "example-topic"
	consumerGroup = "example-consumer-group"
)

// ConsumerGroupHandler 实现 sarama.ConsumerGroupHandler 接口
type ConsumerGroupHandler struct {
	logger *log.Logger
}

// Setup 在新的 session 开始时调用
func (h *ConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	h.logger.Printf("🎯 新会话开始，成员 ID: %s", session.MemberID())
	h.logger.Printf("📋 分配的分区: %v", session.Claims())
	return nil
}

// Cleanup 在 session 结束时调用
func (h *ConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	h.logger.Printf("🔚 会话结束")
	return nil
}

// ConsumeClaim 处理消息
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 注意: 不要在这个函数中启动 goroutine
	// ConsumeClaim 会为每个分区启动一个 goroutine

	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			// 打印消息详情
			h.logger.Printf("📨 收到消息:")
			h.logger.Printf("  Topic: %s", message.Topic)
			h.logger.Printf("  Partition: %d", message.Partition)
			h.logger.Printf("  Offset: %d", message.Offset)
			h.logger.Printf("  Key: %s", string(message.Key))
			h.logger.Printf("  Value: %s", string(message.Value))
			h.logger.Printf("  Timestamp: %s", message.Timestamp.Format("2006-01-02 15:04:05"))

			// 打印消息头
			if len(message.Headers) > 0 {
				h.logger.Printf("  Headers:")
				for _, header := range message.Headers {
					h.logger.Printf("    %s: %s", string(header.Key), string(header.Value))
				}
			}

			// 处理消息
			if err := h.processMessage(message); err != nil {
				h.logger.Printf("❌ 处理消息失败: %v", err)
				// 在生产环境中，可以选择重试或将失败消息发送到死信队列
				continue
			}

			// 标记消息为已处理
			session.MarkMessage(message, "")

			h.logger.Printf("✅ 消息处理完成\n")

		case <-session.Context().Done():
			return nil
		}
	}
}

// processMessage 处理单条消息
func (h *ConsumerGroupHandler) processMessage(message *sarama.ConsumerMessage) error {
	// 这里实现你的业务逻辑
	// 例如: 解析 JSON、写入数据库、调用 API 等

	// 模拟处理
	// time.Sleep(100 * time.Millisecond)

	return nil
}

func main() {
	// 创建日志记录器
	logger := log.New(os.Stdout, "[Consumer] ", log.LstdFlags)

	// 配置消费者
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest                   // 从最新位置开始消费
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky // 使用粘性策略

	// 手动提交 Offset
	config.Consumer.Offsets.AutoCommit.Enable = false

	// 会话超时配置
	config.Consumer.Group.Session.Timeout = 20 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 6 * time.Second

	// Broker 地址列表
	brokers := []string{
		"localhost:19092",
		"localhost:29092",
		"localhost:39092",
	}

	logger.Println("启动消费者...")
	logger.Printf("消费者组: %s", consumerGroup)
	logger.Printf("订阅 Topic: %s", topic)
	logger.Printf("Broker 地址: %v", brokers)

	// 创建消费者组
	consumerGroup, err := sarama.NewConsumerGroup(brokers, consumerGroup, config)
	if err != nil {
		log.Fatalf("无法创建消费者组: %v", err)
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建消费者处理器
	handler := &ConsumerGroupHandler{
		logger: logger,
	}

	// 启动消费者
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// Consume 会一直阻塞，直到发生 rebalance 或 context 取消
			if err := consumerGroup.Consume(ctx, []string{topic}, handler); err != nil {
				logger.Printf("消费错误: %v", err)
			}

			// 检查 context 是否已取消
			if ctx.Err() != nil {
				return
			}

			logger.Println("重新加入消费者组...")
		}
	}()

	// 处理错误
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range consumerGroup.Errors() {
			logger.Printf("❌ 消费者错误: %v", err)
		}
	}()

	logger.Println("✅ 消费者已启动，等待消息...")
	logger.Println("按 Ctrl+C 退出")

	// 等待退出信号
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	logger.Println("收到退出信号，正在关闭...")
	cancel()
	wg.Wait()

	if err := consumerGroup.Close(); err != nil {
		logger.Printf("关闭消费者组失败: %v", err)
	} else {
		logger.Println("消费者已关闭")
	}
}
