package http

import (
	"fmt"

	"github.com/alisstaki/interview-exercise/http"
	apiv1 "github.com/alisstaki/interview-exercise/service/internal/http/internal/api/v1"
	apiv2 "github.com/alisstaki/interview-exercise/service/internal/http/internal/api/v2"
	"github.com/gorilla/mux"
)

type Runner interface {
	apiv1.Runner
	apiv2.Runner
}

// StartAPI creates a new instance of apiv1.Service.
func StartAPI(addr string, runner Runner) error {
	serviceV1 := apiv1.NewService(runner)
	serviceV2 := apiv2.NewService(runner)

	router := mux.NewRouter()
	router.HandleFunc("/run", serviceV1.Run).Methods("POST")
	router.HandleFunc("/v2/run", serviceV2.Run).Methods("POST")
	router.HandleFunc("/status", serviceV1.Status).Methods("GET").Queries("id", "{jobid}")

	fmt.Println("API ready to receive content")
	if err := http.ListenAndServe(addr, router); err != nil {
		return fmt.Errorf("HTTP listenAndServe: %w", err)
	}
	fmt.Println("API stopped")
	return nil
}
