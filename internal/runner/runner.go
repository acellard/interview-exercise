package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	dockerfileinstructions "github.com/moby/buildkit/frontend/dockerfile/instructions"
	dockerfileparser "github.com/moby/buildkit/frontend/dockerfile/parser"
)

type Handler interface {
	ExecJob(fileName string) error
	ReadJobPerformance(jobID string) (float64, error)
}

type handler struct{}

func NewHandler() Handler {
	return &handler{}
}

func (h *handler) ReadJobPerformance(jobID string) (float64, error) {
	// Check job status
	jsonFile, err := os.Open("./data/perf.json")
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
	err = json.Unmarshal([]byte(b), &r)
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

func (h *handler) ExecJob(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	err = h.verifyDockerFile(f)
	if err != nil {
		return err
	}

	err = h.execDockerFile()
	if err != nil {
		return err
	}

	return nil
}

func (h *handler) verifyDockerFile(f *os.File) error {
	result, err := dockerfileparser.Parse(f)
	if err != nil {
		fmt.Println("docker parser :", err)
		return err
	}
	_, _, err = dockerfileinstructions.Parse(result.AST)
	if err != nil {
		fmt.Println("docker instructions :", err)
		return err
	}

	return nil
}

func (h *handler) execDockerFile() error {
	// build docker image
	build := exec.Command("/bin/sh", "-c", "docker build -t input-image .")
	err := build.Run()
	if err != nil {
		fmt.Println("cannot create image based on provided docker file")
		return err
	}

	// once image is created, run docker-compose
	compose := exec.Command("/bin/sh", "-c", "docker-compose up")
	err = compose.Run()
	if err != nil {
		fmt.Println("cannot run container with image based on provided docker file: ", err)
		return err
	}

	return nil
}
