package main

import (
	"fmt"
	"time"
)

func main() {
	channel := make(chan int)
	var sender = func() {
		fmt.Println("data has been send")
		channel <- 1
	}

	var receiver = func() {
		time.Sleep(2000 * time.Millisecond)
		x := <-channel
		fmt.Println("data has been received, ", x)
	}
	go sender()
	go receiver()
	time.Sleep(5 * time.Second)
}
