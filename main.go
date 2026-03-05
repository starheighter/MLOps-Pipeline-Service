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

type Model struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"createdAt"`
	LearningRate float64    `json:"learningRate"`
	Accuracy     float64    `json:"accuracy"`
	Weights      [3]float64 `json:"weights"`
}

type DataSet struct {
	Name   string      `json:"name"`
	Input  [][]float64 `json:"input"`
	Output []float64   `json:"output"`
}

type TrainRequest struct {
	ModelName    string  `json:"modelName"`
	LearningRate float64 `json:"learningRate"`
	DataSet      DataSet `json:"dataSet"`
}

var (
	models = make(map[string]*Model)
	mu     sync.Mutex
)

func createModel() *Model {
	modelId := fmt.Sprintf("%d", rand.Intn(100000))
	model := &Model{
		ID:           modelId,
		Name:         "demo-model",
		Status:       "created",
		CreatedAt:    time.Now(),
		LearningRate: 0.0,
		Accuracy:     0.0,
		Weights:      [3]float64{0.0, 0.0, 0.0},
	}
	models[modelId] = model
	return model
}

func handleTrain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req TrainRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	mu.Lock()
	model := createModel()
	mu.Unlock()
	modelId := r.URL.Path[len("/train/"):]
	mu.Lock()
	modelId = model.ID
	model, ok := models[modelId]
	if !ok {
		mu.Unlock()
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}
	model.Status = "trained"
	mu.Unlock()
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model)
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	modelId := r.URL.Path[len("/status/"):]
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

func handleModels(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	var list []*Model
	for _, m := range models {
		list = append(list, m)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	if os.Getenv("SEED_JOBS") == "true" {
		for i := 0; i < 5; i++ {
			createModel()
		}
	}
	http.HandleFunc("/train", handleTrain)
	http.HandleFunc("/deploy/", handleDeploy)
	http.HandleFunc("/status/", handleStatus)
	http.HandleFunc("/models", handleModels)
	http.HandleFunc("/health", handleHealth)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
