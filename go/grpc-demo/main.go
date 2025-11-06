package main

import (
	"log"
	"os"

	"github.com/MorseWayne/grpc-demo/internal/client"
	"github.com/MorseWayne/grpc-demo/internal/server"
)

func main() {
	if len(os.Args) <= 1 {
		return
	}
	if os.Args[1] == "client" {
		if err := client.Run("localhost:12345"); err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := server.Run("localhost:12345"); err != nil {
		log.Fatal(err)
	}
}
