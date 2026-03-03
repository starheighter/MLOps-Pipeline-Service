package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

type JobStatus string

const (
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

type Job struct {
	ID     string    `json:"id"`
	Status JobStatus `json:"status"`
}

type Model struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var (
	jobs   = make(map[string]*Job)
	models = make(map[string]*Model)
	mu     sync.Mutex
)

func createJob() *Job {
	jobId := fmt.Sprintf("%d", rand.Intn(100000))
	job := &Job{ID: jobId, Status: StatusRunning}
	mu.Lock()
	jobs[jobId] = job
	mu.Unlock()
	go func() {
		time.Sleep(30 * time.Second)
		mu.Lock()
		job.Status = StatusCompleted
		models[jobId] = &Model{
			ID:          jobId,
			Description: "Simuliertes Modell",
			Status:      "trained",
		}
		mu.Unlock()
	}()
	fmt.Println(job)
	return job
}

func handleTrain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	job := createJob()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	jobId := r.URL.Path[len("/status/"):]
	mu.Lock()
	job, ok := jobs[jobId]
	mu.Unlock()
	if !ok {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func handleModel(w http.ResponseWriter, r *http.Request) {
	modelId := r.URL.Path[len("/model/"):]
	mu.Lock()
	model, ok := models[modelId]
	mu.Unlock()
	if !ok {
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model)
}

func handleDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	modelId := r.URL.Path[len("/deploy/"):]
	mu.Lock()
	model, ok := models[modelId]
	if !ok {
		mu.Unlock()
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}
	model.Status = "deployed"
	mu.Unlock()
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Model deployed successfully",
		"modelId":   model.ID,
		"newStatus": model.Status,
	})
}

func main() {
	if os.Getenv("SEED_JOBS") == "true" {
		for i := 0; i < 5; i++ {
			createJob()
		}
	}
	http.HandleFunc("/train", handleTrain)
	http.HandleFunc("/status/", handleStatus)
	http.HandleFunc("/model/", handleModel)
	http.HandleFunc("/deploy/", handleDeploy)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
