package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/alisstaki/interview-exercise/internal/runner"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Service has a router
type Service struct {
	Router        *mux.Router
	RunnerHandler runner.Handler
	jobStatuses   map[string]string
}

// New creates a new instance of Service
func New(router *mux.Router, runnerHandler runner.Handler) *Service {
	return &Service{
		Router:        router,
		RunnerHandler: runnerHandler,
		jobStatuses:   map[string]string{},
	}
}

// DefineHandlers defines all routes handled by Service
func (s *Service) DefineHandlers() {
	s.Router.HandleFunc("/run", s.Run).Methods("POST")
	s.Router.HandleFunc("/status", s.Status).Methods("GET").Queries("id", "{jobid}")
}

func (s *Service) Run(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	file := r.FormValue("file")

	err := ioutil.WriteFile("dockerfile", []byte(file), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	id := uuid.New()
	go func(jobID uuid.UUID) {
		err := s.RunnerHandler.ExecJob("dockerfile")
		if err != nil {
			fmt.Println(err)
			s.jobStatuses[jobID.String()] = "Failed"
		} else {
			s.jobStatuses[jobID.String()] = "Successful"
		}
	}(id)

	w.Write([]byte(fmt.Sprintf("Job was launched with id %v", id.String())))
}

func (s *Service) Status(w http.ResponseWriter, r *http.Request) {
	// Get request query parameters
	vars := mux.Vars(r)
	jobid := vars["jobid"]

	// check status
	var status string
	if _, ok := s.jobStatuses[jobid]; !ok {
		http.Error(w, fmt.Sprintf("Job with id %s was not found", jobid), http.StatusNotFound)
	}
	status = s.jobStatuses[jobid]

	perf, err := s.RunnerHandler.ReadJobPerformance(jobid)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(fmt.Sprintf("Job %s was %s - Model performance : %f", jobid, status, perf)))
}
