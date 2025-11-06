package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/morsewayne/kafka-demo/examples/08-order-processing/models"
)

const orderTopic = "order-events"

func main() {
	logger := log.New(os.Stdout, "[OrderService] ", log.LstdFlags)

	// 配置生产者
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	brokers := []string{"localhost:19092", "localhost:29092", "localhost:39092"}

	logger.Println("🚀 启动订单服务...")
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("创建生产者失败: %v", err)
	}
	defer producer.Close()

	logger.Println("✅ 订单服务已启动")

	// 模拟创建订单
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		orderNum := 1
		for range ticker.C {
			if err := createOrder(producer, logger, orderNum); err != nil {
				logger.Printf("❌ 创建订单失败: %v", err)
			}
			orderNum++
		}
	}()

	// 等待退出信号
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	logger.Println("收到退出信号，关闭服务...")
}

func createOrder(producer sarama.SyncProducer, logger *log.Logger, orderNum int) error {
	// 生成订单
	traceID := uuid.New().String()
	orderID := fmt.Sprintf("ORD-%06d", orderNum)
	userID := fmt.Sprintf("USER-%03d", orderNum%10)

	order := models.OrderCreated{
		EventType: models.EventOrderCreated,
		OrderID:   orderID,
		UserID:    userID,
		Items: []models.OrderItem{
			{
				ProductID: "PROD-001",
				Name:      "iPhone 15 Pro",
				Quantity:  1,
				Price:     999.99,
			},
			{
				ProductID: "PROD-002",
				Name:      "AirPods Pro",
				Quantity:  1,
				Price:     249.99,
			},
		},
		TotalAmount: 1249.98,
		Timestamp:   time.Now(),
		TraceID:     traceID,
	}

	// 序列化
	jsonData, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("序列化订单失败: %w", err)
	}

	// 发送消息
	msg := &sarama.ProducerMessage{
		Topic: orderTopic,
		Key:   sarama.StringEncoder(orderID), // 使用订单 ID 作为 Key，确保有序
		Value: sarama.ByteEncoder(jsonData),
		Headers: []sarama.RecordHeader{
			{Key: []byte("event_type"), Value: []byte(models.EventOrderCreated)},
			{Key: []byte("trace_id"), Value: []byte(traceID)},
		},
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}

	logger.Printf("📦 订单已创建: OrderID=%s, UserID=%s, Amount=%.2f, TraceID=%s",
		orderID, userID, order.TotalAmount, traceID)
	logger.Printf("   发送到 Partition=%d, Offset=%d\n", partition, offset)

	return nil
}
