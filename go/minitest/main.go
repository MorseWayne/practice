package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func producer(factor int, out chan<- int) {
	for i := 0; ; i++ {
		out <- i * factor
	}
}

func consumer(in <-chan int) {
	for num := range in {
		println("get num: ", num)
	}
}

func main() {
	channel := make(chan int)
	go producer(3, channel)
	go producer(5, channel)
	go consumer(channel)

	errChan := make(chan os.Signal, 1)
	signal.Notify(errChan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <-errChan)
}
