package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	fmt.Println("Welcome to the g-commander. Orthodox file manager")
	fmt.Println("Press Ctrl+C to exit")

	<-ctx.Done()
}
