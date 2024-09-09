package main

import (
	"log"
	"github.com/ulfaric/gomodbus/socket"
)

func main() {
	w := socket.NewSocket(1993)
	err := w.Start()
	if err != nil {
		log.Fatalf("Failed to start wrapper: %v", err)
	}
}