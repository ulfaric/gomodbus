package main

import (
	"log"
	"github.com/ulfaric/gomodbus/socket"
)

func main() {
	s := socket.NewSocket(1993)
	err := s.Start()
	if err != nil {
		log.Fatalf("Failed to start wrapper: %v", err)
	}
}