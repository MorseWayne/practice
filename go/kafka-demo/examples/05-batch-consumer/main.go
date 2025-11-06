package batchconsumer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

const (
	topic         = "batch-processing-topic"
	consumerGroup = "batch-consumer-group"
	batchSize     = 50 // 批处理大小
)

type OrderEvent struct {
	OrderID    string    `json:"order_id"`
	UserID     string    `json:"user_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	CreateTime time.Time `json:"create_time"`
}

type BatchConsumerHandler struct {
	logger *log.Logger
}

func (h *BatchConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	h.logger.Printf("🎯 新会话开始")
	return nil
}

func (h *BatchConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	h.logger.Printf("🔚 会话结束")
	return nil
}

func (h *BatchConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	batch := make([]*sarama.ConsumerMessage, 0, batchSize)
	ticker := time.NewTicker(5 * time.Second) // 最多等待 5 秒
	defer ticker.Stop()

	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			batch = append(batch, message)

			// 达到批次大小，立即处理
			if len(batch) >= batchSize {
				h.processBatch(session, batch)
				batch = batch[:0]
				ticker.Reset(5 * time.Second)
			}

		case <-ticker.C:
			// 超时，处理当前批次（即使未满）
			if len(batch) > 0 {
				h.logger.Printf("⏰ 批次超时，处理 %d 条消息", len(batch))
				h.processBatch(session, batch)
				batch = batch[:0]
			}

		case <-session.Context().Done():
			// 会话结束前处理剩余消息
			if len(batch) > 0 {
				h.logger.Printf("🔄 会话结束，处理剩余 %d 条消息", len(batch))
				h.processBatch(session, batch)
			}
			return nil
		}
	}
}

func (h *BatchConsumerHandler) processBatch(session sarama.ConsumerGroupSession, batch []*sarama.ConsumerMessage) {
	startTime := time.Now()

	h.logger.Printf("📦 开始处理批次，消息数: %d", len(batch))

	// 解析所有消息
	orders := make([]OrderEvent, 0, len(batch))
	for _, msg := range batch {
		var order OrderEvent
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			h.logger.Printf("❌ 解析消息失败: %v", err)
			continue
		}
		orders = append(orders, order)
	}

	// 批量处理业务逻辑
	// 例如: 批量插入数据库、批量调用 API 等
	if err := h.batchProcess(orders); err != nil {
		h.logger.Printf("❌ 批量处理失败: %v", err)
		return
	}

	// 标记最后一条消息
	lastMsg := batch[len(batch)-1]
	session.MarkMessage(lastMsg, "")

	duration := time.Since(startTime)
	h.logger.Printf("✅ 批次处理完成，耗时: %v, 速率: %.2f 条/秒\n",
		duration, float64(len(batch))/duration.Seconds())
}

func (h *BatchConsumerHandler) batchProcess(orders []OrderEvent) error {
	// 模拟批量处理
	// 例如: 批量写入数据库
	time.Sleep(100 * time.Millisecond)

	// 统计信息
	totalAmount := 0.0
	userMap := make(map[string]int)

	for _, order := range orders {
		totalAmount += order.Amount
		userMap[order.UserID]++
	}

	h.logger.Printf("  订单数: %d, 总金额: %.2f, 用户数: %d",
		len(orders), totalAmount, len(userMap))

	return nil
}

func main() {
	logger := log.New(os.Stdout, "[BatchConsumer] ", log.LstdFlags)

	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // 从头开始
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	config.Consumer.Offsets.AutoCommit.Enable = false // 手动提交

	// 增加拉取大小，适合批处理
	config.Consumer.Fetch.Default = 1024 * 1024 // 1MB
	config.Consumer.MaxProcessingTime = 30 * time.Second

	brokers := []string{"localhost:19092", "localhost:29092", "localhost:39092"}

	logger.Println("启动批量消费者...")
	logger.Printf("批处理大小: %d", batchSize)

	consumerGroup, err := sarama.NewConsumerGroup(brokers, consumerGroup, config)
	if err != nil {
		log.Fatalf("创建消费者组失败: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := &BatchConsumerHandler{logger: logger}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, handler); err != nil {
				logger.Printf("消费错误: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range consumerGroup.Errors() {
			logger.Printf("❌ 消费者错误: %v", err)
		}
	}()

	logger.Println("✅ 批量消费者已启动")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	logger.Println("收到退出信号...")
	cancel()
	wg.Wait()
	consumerGroup.Close()
	logger.Println("批量消费者已关闭")
}
