package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/xtsank/mypills-super-service/src/internal/benchmark/lab3"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return lab3.Run(ctx, lab3.Config{
		DSN: os.Getenv("LAB3_DSN"), RunID: os.Getenv("LAB3_RUN_ID"),
		Profile: os.Getenv("LAB3_PROFILE"), ResultsDir: "/results",
		CPU: os.Getenv("LAB3_CPU"),
	})
}
