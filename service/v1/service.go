package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/alisstaki/interview-exercise/internal/content"
	"github.com/gorilla/mux"
)

// Service has a router
type Service struct {
	Router         *mux.Router
	ContentHandler content.Handler
}

// New creates a new instance of Service
func New(router *mux.Router, contentHandler content.Handler) *Service {
	return &Service{
		Router:         router,
		ContentHandler: contentHandler,
	}
}

// DefineHandlers defines all routes handled by Service
func (s *Service) DefineHandlers() {
	s.Router.HandleFunc("/content", s.handleContent).Methods("POST")
}

func (s *Service) handleContent(w http.ResponseWriter, r *http.Request) {
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

	// call
	perf, err := s.ContentHandler.Handle("dockerfile")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(fmt.Sprintf("%f", perf)))
}
