package main

import (
	"context"
	"os"
	"os/signal"
)

func GracefulShutdown(fn func(), sigs ...os.Signal) <-chan struct{} {
	stop := make(chan struct{})
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, sigs...)

	go func() {
		<-sigChan

		signal.Stop(sigChan)

		fn()

		close(sigChan)
		close(stop)
	}()

	return stop
}

func GracefulShutdownCtx(ctx context.Context, fn func(), sigs ...os.Signal) (context.Context, context.CancelFunc) {
	ctx2, cancel := signal.NotifyContext(ctx, sigs...)

	go func() {
		<-ctx2.Done()
		fn()
	}()

	return ctx2, cancel
}
