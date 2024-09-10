package main

import (
	"github.com/ulfaric/gomodbus/socket"
)

func main() {
	socket := socket.NewSocket()
	socket.Start()
}