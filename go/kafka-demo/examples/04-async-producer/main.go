package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

const topic = "batch-processing-topic"

// OrderEvent 订单事件
type OrderEvent struct {
	OrderID    string    `json:"order_id"`
	UserID     string    `json:"user_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	CreateTime time.Time `json:"create_time"`
}

func main() {
	logger := log.New(os.Stdout, "[BatchProducer] ", log.LstdFlags)

	// 配置异步生产者
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForLocal // 只等待 Leader 确认（更快）
	config.Producer.Compression = sarama.CompressionSnappy

	// 批处理配置
	config.Producer.Flush.Messages = 100                     // 批量发送 100 条消息
	config.Producer.Flush.Frequency = 100 * time.Millisecond // 或每 100ms 发送一次
	config.Producer.Flush.MaxMessages = 1000                 // 最大批次大小

	brokers := []string{"localhost:19092", "localhost:29092", "localhost:39092"}

	logger.Println("创建异步生产者...")
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("创建生产者失败: %v", err)
	}
	defer producer.Close()

	logger.Println("✅ 异步生产者已启动")

	// 统计信息
	var (
		successCount int64
		errorCount   int64
		mu           sync.Mutex
	)

	// 处理成功的消息
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			mu.Lock()
			successCount++
			count := successCount
			mu.Unlock()

			if count%100 == 0 {
				logger.Printf("✅ 已成功发送 %d 条消息", count)
			}
		}
	}()

	// 处理错误
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			mu.Lock()
			errorCount++
			mu.Unlock()
			logger.Printf("❌ 发送失败: %v", err.Err)
		}
	}()

	// 发送消息
	logger.Println("开始批量发送消息...")
	startTime := time.Now()

	// 模拟批量生成订单
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 1; i <= 1000; i++ {
			order := OrderEvent{
				OrderID:    fmt.Sprintf("ORDER-%06d", i),
				UserID:     fmt.Sprintf("USER-%04d", i%100),
				Amount:     float64(i) * 10.5,
				Status:     "CREATED",
				CreateTime: time.Now(),
			}

			jsonData, _ := json.Marshal(order)

			msg := &sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.StringEncoder(order.UserID), // 按用户 ID 分区
				Value: sarama.ByteEncoder(jsonData),
			}

			// 异步发送（非阻塞）
			producer.Input() <- msg

			// 每 100 条消息打印进度
			if i%100 == 0 {
				logger.Printf("📤 已投递 %d 条消息到发送队列", i)
			}
		}

		logger.Println("所有消息已投递到发送队列")
	}()

	// 等待退出信号
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		logger.Println("收到退出信号，等待消息发送完成...")
		producer.AsyncClose()
	}()

	wg.Wait()

	duration := time.Since(startTime)
	logger.Printf("\n📊 发送统计:")
	logger.Printf("  总耗时: %v", duration)
	logger.Printf("  成功: %d 条", successCount)
	logger.Printf("  失败: %d 条", errorCount)
	logger.Printf("  速率: %.2f 条/秒", float64(successCount)/duration.Seconds())
}
