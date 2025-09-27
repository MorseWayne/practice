package goroutine

type PrimeSieve struct {
}

func NewPrimeSieve() *PrimeSieve {
	return &PrimeSieve{}
}

// 指定数量生成素数
func (ps *PrimeSieve) Generate(want int) []int {
	res := make([]int, 0)
	ch := ps.generateNatural() // 自然数序列: 2, 3, 4, ...
	for range want {
		prime := <-ch // 新出现的素数
		res = append(res, prime)
		ch = ps.filter(ch, prime) // 基于新素数构造的过滤器
	}
	return res
}

func (ps *PrimeSieve) generateNatural() chan int {
	channel := make(chan int)
	go func() {
		for i := 2; i < 100; i++ {
			channel <- i
		}
	}()
	return channel
}

func (ps *PrimeSieve) filter(in <-chan int, prime int) chan int {
	out := make(chan int)
	go func() {
		for i := range in {
			if i%prime != 0 {
				out <- i
			}
		}
	}()
	return out
}
