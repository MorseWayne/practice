package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("✓ connected to nats server")

	// 示例1：简单订阅
	subject := "demo.sample"
	fmt.Printf("✓ start subscribe[simple]: %s\n", subject)
	sub, err := conn.Subscribe(subject, func(msg *nats.Msg) {
		fmt.Printf("✓ receive[simple]: %s, %s\n", msg.Subject, string(msg.Data))
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe()

	// 示例2：通过通配符订阅
	fmt.Printf("✓ start subscribe[wildcard]: demo.*\n")
	wildcardSub, err := conn.Subscribe("demo.*", func(msg *nats.Msg) {
		fmt.Printf("✓ receive[wildcard]: %s, %s\n", msg.Subject, string(msg.Data))
	})
	if err != nil {
		log.Fatal(err)
	}
	defer wildcardSub.Unsubscribe()

	// 示例3：通过队列组订阅
	fmt.Printf("✓ start subscribe[queue]: \n")
	queueSub, err := conn.QueueSubscribe("demo.queue", "workers", func(msg *nats.Msg) {
		fmt.Printf("✓ receive[queue]: %s : %s\n", msg.Subject, string(msg.Data))
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("✓ msg process finished\n")
	})
	if err != nil {
		log.Fatal(err)
	}
	defer queueSub.Unsubscribe()
	fmt.Println("✓ start waitting msg, ctrl + c exit!")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Printf("clearing...\n")
}
