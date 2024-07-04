package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	server := NewServer("localhost:12345", &DummyKV{})
	server.Start()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	server.Close()
	fmt.Println("terminated")
}
