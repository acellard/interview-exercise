package main

import (
	"fmt"
	"time"

	"github.com/alisstaki/interview-exercise/service/internal/filesystem"
	"github.com/alisstaki/interview-exercise/service/internal/http"
)

func main() {
	runner := filesystem.NewRunner(10)
	// Ensure jobs running in background are finishing properly before exit.
	defer runner.Shutdown(10 * time.Second)

	// Listen to port 8080 to handle requests
	if err := http.StartAPI(":8080", runner); err != nil {
		panic(fmt.Errorf("HTTP server exited, error: %w", err))
	}
}
