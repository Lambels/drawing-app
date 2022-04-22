package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	server := NewServer(":8080")
	go server.Open()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGQUIT)

	<-c

	server.Close()
}
