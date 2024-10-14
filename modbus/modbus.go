package main

// import (
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"github.com/ulfaric/gomodbus/socket"
// )

// func main() {
// 	socket := socket.NewSocket()
// 	socket.Start()

// 	// Create a channel to listen for termination signals
// 	sigChan := make(chan os.Signal, 1)
// 	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

// 	// Block until a signal is received
// 	<-sigChan

// 	// Gracefully shut down the socket
// 	socket.Stop()
// }