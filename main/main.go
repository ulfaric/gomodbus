package main

import (
	"log"
	"github.com/ulfaric/gomodbus/wrapper"
)

func main() {
	w := wrapper.NewWrapper(1993)
	err := w.Start()
	if err != nil {
		log.Fatalf("Failed to start wrapper: %v", err)
	}
}