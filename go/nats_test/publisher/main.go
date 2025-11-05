package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	println("✓ connected to nats server\n")
	println("✓ ready to publish msg\n")

	// 1. 简单发布
	subject := "demo.sample"
	println("✓ ready to publish in simple way\n")
	for i := range 3 {
		msg := fmt.Sprintf("Hello nats #%d", i)
		err := nc.Publish(subject, []byte(msg))
		if err != nil {
			log.Printf("	× publish failed: %v\n", err)
			continue
		}
		fmt.Printf("	✓ published: %s\n", msg)
		time.Sleep(500 * time.Millisecond)
	}

	// 2. 发布多个主题
	println("✓ ready to publish multi subjects")
	subjects := []string{"demo.foo", "demo.bar", "demo.baz"}
	for _, subject := range subjects {
		msg := fmt.Sprintf("this is a message to %s", subject)
		nc.Publish(subject, []byte(msg))
		fmt.Printf("	✓ published %s : %s\n", subject, msg)
		time.Sleep(500 * time.Millisecond)
	}

	// 3. 发布到队列
	println("✓ ready to publish to queue workers\n")
	queueSubject := "demo.queue"
	for i := range 5 {
		msg := fmt.Sprintf("Queue task #%d", i)
		nc.Publish(queueSubject, []byte(msg))
		fmt.Printf("	✓ published %s : %s\n", subject, msg)
		time.Sleep(300 * time.Millisecond)
	}

	nc.Flush()
	println("all message has been published!!!")
}
