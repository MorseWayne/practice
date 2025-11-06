package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

const (
	topic = "example-topic"
)

func main() {
	// 创建日志记录器
	logger := log.New(os.Stdout, "[Producer] ", log.LstdFlags)

	// 配置 Kafka 生产者
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Producer.Return.Successes = true                // 等待服务器确认
	config.Producer.Return.Errors = true                   // 返回错误
	config.Producer.RequiredAcks = sarama.WaitForAll       // 等待所有副本确认（最高可靠性）
	config.Producer.Retry.Max = 3                          // 失败重试次数
	config.Producer.Compression = sarama.CompressionSnappy // 使用 Snappy 压缩

	// Broker 地址列表
	brokers := []string{
		"localhost:19092",
		"localhost:29092",
		"localhost:39092",
	}

	logger.Printf("连接到 Kafka 集群...")
	logger.Printf("Broker 地址: %v", brokers)

	// 创建同步生产者
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("无法创建生产者: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			logger.Printf("关闭生产者失败: %v", err)
		} else {
			logger.Println("关闭生产者")
		}
	}()

	logger.Println("成功连接到 Kafka")

	// 设置信号处理，优雅退出
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// 发送消息
	logger.Println("开始发送消息...")

	go func() {
		messageCount := 0
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				messageCount++

				// 构造消息
				message := fmt.Sprintf("消息 #%d - 时间: %s", messageCount, time.Now().Format(time.RFC3339))

				// 创建 Kafka 消息
				msg := &sarama.ProducerMessage{
					Topic: topic,
					Key:   sarama.StringEncoder(fmt.Sprintf("key-%d", messageCount%3)), // 使用 Key 确保相同 Key 的消息到同一分区
					Value: sarama.StringEncoder(message),
					Headers: []sarama.RecordHeader{
						{
							Key:   []byte("source"),
							Value: []byte("simple-producer"),
						},
						{
							Key:   []byte("version"),
							Value: []byte("1.0"),
						},
					},
					Timestamp: time.Now(),
				}

				// 发送消息（同步）
				partition, offset, err := producer.SendMessage(msg)
				if err != nil {
					logger.Printf("❌ 发送消息失败: %v", err)
				} else {
					logger.Printf("✅ 消息已发送 -> Topic: %s, Partition: %d, Offset: %d, Key: key-%d",
						topic, partition, offset, messageCount%3)
				}

				// 发送 10 条消息后停止
				if messageCount >= 10 {
					logger.Printf("总共发送了 %d 条消息", messageCount)
					signals <- syscall.SIGTERM
					return
				}

			case <-signals:
				return
			}
		}
	}()

	// 等待退出信号
	<-signals
	logger.Println("收到退出信号，正在关闭...")
}
