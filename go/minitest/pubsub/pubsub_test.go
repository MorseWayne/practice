package pubsub_test

import (
	"fmt"
	"minitest/pubsub"
	"strings"
	"testing"
	"time"
)

func TestPubsub(t *testing.T) {
	pub := pubsub.NewPublisher(time.Second, 10)
	defer pub.Close()

	sub1 := pub.Subscrible()
	sub2 := pub.SubscribleTopic(func(msg interface{}) bool {
		if s, ok := msg.(string); ok {
			return strings.Contains(s, "key")
		}
		return false
	})

	pub.Publish("hello world!")
	pub.Publish("hello key!")

	go func() {
		for msg := range sub1 {
			fmt.Println("sub1: ", msg)
		}
	}()

	go func() {
		for msg := range sub2 {
			fmt.Println("sub2: ", msg)
		}
	}()

	time.Sleep(time.Second * 3)
}
