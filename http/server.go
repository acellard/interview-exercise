package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func ListenAndServe(addr string, handler http.Handler) (err error) {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Handle graceful shutdown
	waitShutdown := make(chan struct{})
	ctx, cancelWaitShutdown := context.WithCancel(context.Background())
	defer func() {
		cancelWaitShutdown()
		<-waitShutdown
	}()
	go func() {
		defer close(waitShutdown)

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		select {
		case <-c:
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second) // use a shutdown timeout
			defer cancel()

			fmt.Println("HTTP server is shutting down after receiving an OS interruption signal")
			if err := server.Shutdown(ctx); err != nil {
				fmt.Println("HTTP server shutdown, error:", err.Error())
			} else {
				fmt.Println("HTTP server has gracefully shutdown")
			}
		case <-ctx.Done():
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		switch err {
		case http.ErrServerClosed:
		default:
			return fmt.Errorf("HTTP server exited: %w", err)
		}
	}
	return nil
}
