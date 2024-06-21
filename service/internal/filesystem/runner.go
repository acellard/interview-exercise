package filesystem

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	dockerfileinstructions "github.com/moby/buildkit/frontend/dockerfile/instructions"
	dockerfilelinter "github.com/moby/buildkit/frontend/dockerfile/linter"
	dockerfileparser "github.com/moby/buildkit/frontend/dockerfile/parser"
	"golang.org/x/sync/semaphore"
)

type Runner struct {
	maxJobHandle int32
	jobLimiter   *semaphore.Weighted

	jobStatusLock sync.RWMutex
	jobStatuses   map[string]string

	jobId int32
}

func NewRunner(maxJobHandle int32) *Runner {
	return &Runner{
		maxJobHandle: maxJobHandle,
		jobLimiter:   semaphore.NewWeighted(int64(maxJobHandle)),
		jobStatuses:  map[string]string{},
	}
}

func (r *Runner) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := r.jobLimiter.Acquire(ctx, int64(r.maxJobHandle)); err != nil {
		return fmt.Errorf("waiting runner stop: %w", err)
	}
	fmt.Println("Runner job list clear")
	return nil
}

func (h *Runner) ReadJobPerformance(jobID string) (float64, error) {
	// Check job status
	jsonFile, err := os.Open("./data/perf-" + jobID + ".json")
	if err != nil {
		fmt.Println("cannot read perf: ", err)
		return 0, err
	}
	defer jsonFile.Close()

	b, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("cannot read json file: ", err)
		return 0, err
	}
	var r map[string]interface{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		fmt.Println("cannot unmarshall json file: ", err)
		return 0, err
	}

	switch perf := r["perf"].(type) {
	case float64:
		return perf, nil
	default:
		return 0, fmt.Errorf("unknown performance type")
	}
}

func (s *Runner) GetJobStatus(jobID string) (string, bool) {
	s.jobStatusLock.RLock()
	defer s.jobStatusLock.RUnlock()
	status, exist := s.jobStatuses[jobID]
	return status, exist
}

func (s *Runner) setJobStatus(jobID, status string) {
	s.jobStatusLock.Lock()
	defer s.jobStatusLock.Unlock()
	s.jobStatuses[jobID] = status
}

func (h *Runner) ExecJob(ctx context.Context, content []byte) (string, error) {
	newJobID := atomic.AddInt32(&h.jobId, 1)
	jobID := strconv.Itoa(int(newJobID))
	if err := h.jobLimiter.Acquire(ctx, 1); err != nil {
		return "", fmt.Errorf("runner exec new job aquire semaphore: %w", err)
	}
	go func() {
		defer h.jobLimiter.Release(1)
		if err := h.execJob(jobID, content); err != nil {
			h.setJobStatus(jobID, "Failed")
			return
		}
		h.setJobStatus(jobID, "Successful")
	}()
	return jobID, nil
}

func (h *Runner) execJob(jobID string, content []byte) error {

	if err := h.verifyDockerFile(content); err != nil {
		return err
	}

	if err := h.execDockerFile(jobID, content); err != nil {
		return err
	}

	return nil
}

func (h *Runner) verifyDockerFile(content []byte) error {
	result, err := dockerfileparser.Parse(bytes.NewBuffer(content))
	if err != nil {
		fmt.Println("docker parser :", err)
		return err
	}
	_, _, err = dockerfileinstructions.Parse(result.AST, &dockerfilelinter.Linter{})
	if err != nil {
		fmt.Println("docker instructions :", err)
		return err
	}

	return nil
}

func (h *Runner) execDockerFile(jobID string, content []byte) error {
	filename := "dockerfile-" + jobID

	err := os.WriteFile(filename, content, os.ModePerm)
	if err != nil {
		return fmt.Errorf("write file")
	}
	defer os.Remove(filename)

	// build docker image
	build := exec.Command("/bin/sh", "-c", "docker build --build-arg version="+jobID+" -t input-image:"+jobID+" -f "+filename+" .")
	err = build.Run()
	if err != nil {
		fmt.Println("cannot create image based on provided docker file:", err.Error(), "jobID:", jobID)
		return err
	}
	defer func() {
		removeImage := exec.Command("/bin/sh", "-c", "docker rmi input-image:"+jobID)
		err = removeImage.Run()
		if err != nil {
			fmt.Println("cannot remove image based on provided docker file:", err.Error())
		}
	}()

	// once image is created, run docker-compose
	compose := exec.Command("/bin/sh", "-c", "COMPOSE_PROJECT_NAME=runner-"+jobID+" VERSION="+jobID+" docker compose up")
	err = compose.Run()
	if err != nil {
		fmt.Println("cannot run container with image based on provided docker file: VERSION="+jobID+" docker compose up:", err.Error())
		return err
	}
	// once image is created, run docker-compose
	rm := exec.Command("/bin/sh", "-c", "COMPOSE_PROJECT_NAME=runner-"+jobID+" VERSION="+jobID+" docker compose down")
	err = rm.Run()
	if err != nil {
		fmt.Println("cannot run container with image based on provided docker file: VERSION="+jobID+" docker compose up:", err.Error())
		return err
	}

	return nil
}
