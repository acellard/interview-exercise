package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/alisstaki/interview-exercise/internal/runner"
	service "github.com/alisstaki/interview-exercise/service/v1"
)

func main() {
	s := service.New(mux.NewRouter(), runner.NewHandler())

	s.DefineHandlers()

	fmt.Println("API ready to receive content...")

	// Listen to port 8080 to handle requests
	log.Fatal(http.ListenAndServe(":8080", s.Router))
}
