package goroutine

import (
	"context"
	"sync"
)

type PrimeSieve struct {
	wg *sync.WaitGroup
}

func NewPrimeSieve() *PrimeSieve {
	return &PrimeSieve{
		wg: &sync.WaitGroup{},
	}

}

// 指定数量生成素数
func (ps *PrimeSieve) Generate(want int) []int {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	res := make([]int, 0)
	ch := ps.generateNatural(ctx, &wg) // 自然数序列: 2, 3, 4, ...
	for i := 0; i < want; i++ {
		prime, ok := <-ch // 新出现的素数，检查channel是否关闭
		if !ok {
			// channel已关闭，没有更多素数了
			break
		}
		res = append(res, prime)
		ch = ps.filter(ctx, ch, prime, &wg) // 基于新素数构造的过滤器
	}
	cancel()
	wg.Wait()
	return res
}

func (ps *PrimeSieve) generateNatural(ctx context.Context, wg *sync.WaitGroup) chan int {
	channel := make(chan int)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(channel)
		for i := 2; i < 100; i++ {
			select {
			case <-ctx.Done():
				return
			case channel <- i:
			}
		}
	}()
	return channel
}

func (ps *PrimeSieve) filter(ctx context.Context, in <-chan int, prime int, wg *sync.WaitGroup) chan int {
	out := make(chan int)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(out)
		for i := range in {
			if i%prime != 0 {
				select {
				case <-ctx.Done():
					return
				case out <- i:

				}
			}

		}
	}()
	return out
}
