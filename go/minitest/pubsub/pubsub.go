package pubsub

import (
	"sync"
	"time"
)

type (
	Subscribler chan interface{}
	TopicFunc   func(v interface{}) bool
)

type Publisher struct {
	RwMtx        sync.RWMutex
	BufferSize   int
	Timeout      time.Duration
	Subscriblers map[Subscribler]TopicFunc
}

// 创建发布者
func NewPublisher(timeout time.Duration, bufferSize int) *Publisher {
	return &Publisher{
		BufferSize:   bufferSize,
		Timeout:      timeout,
		Subscriblers: make(map[Subscribler]TopicFunc),
	}
}

func (p *Publisher) Publish(message interface{}) {
	p.RwMtx.Lock()
	defer p.RwMtx.Unlock()

	var wg sync.WaitGroup
	for sub, topic := range p.Subscriblers {
		wg.Add(1)
		go p.sendTopic(sub, topic, message, &wg)
	}
	wg.Wait()
}

// 带topic过滤的订阅
func (p *Publisher) SubscribleTopic(topic TopicFunc) Subscribler {
	channel := make(Subscribler, p.BufferSize)
	p.RwMtx.Lock()
	defer p.RwMtx.Unlock()
	p.Subscriblers[channel] = topic
	return channel
}

// 创建订阅
func (p *Publisher) Subscrible() Subscribler {
	return p.SubscribleTopic(nil)
}

// 取消订阅
func (p *Publisher) Evict(sub Subscribler) {
	p.RwMtx.Lock()
	defer p.RwMtx.Unlock()
	delete(p.Subscriblers, sub)
	close(sub)
}

func (p *Publisher) Close() {
	p.RwMtx.Lock()
	defer p.RwMtx.Unlock()
	for sub := range p.Subscriblers {
		delete(p.Subscriblers, sub)
		close(sub)
	}
}

// 带topic发布消息
func (p *Publisher) sendTopic(sub Subscribler, topic TopicFunc, message interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	if topic != nil && !topic(message) {
		return
	}
	select {
	case sub <- message:
	case <-time.After(p.Timeout):
	}
}
