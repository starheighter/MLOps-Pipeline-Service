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
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	LearningRate float64   `json:"learningRate"`
	Loss         float64   `json:"loss"`
	Weights      []float64 `json:"weights"`
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
		LearningRate: 0.05,
		Loss:         0.0,
		Weights:      []float64{0.0, 0.0},
	}
	mu.Lock()
	models[modelId] = model
	mu.Unlock()
	return model
}

func handleTrain(w http.ResponseWriter, r *http.Request) {
	modelId := r.URL.Path[len("/train/"):]
	trainingData := &DataSet{
		Name:   "demo-dataset",
		Input:  [][]float64{{0.2, 0.4}, {0.1, 0.8}, {0.3, 0.3}, {0.5, 0.2}},
		Output: []float64{1.0, 1.1, 1.2, 1.7},
	}
	mu.Lock()
	model, ok := models[modelId]
	if !ok {
		mu.Unlock()
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}
	for trainLoop := 0; trainLoop < 50; trainLoop++ {
		for exampleIndex := 0; exampleIndex < len(trainingData.Input); exampleIndex++ {
			weightedSum := 0.0
			loss := 0.0
			for featureIndex := 0; featureIndex < len(trainingData.Input[exampleIndex]); featureIndex++ {
				weightedSum += model.Weights[featureIndex] * trainingData.Input[exampleIndex][featureIndex]
			}
			loss += weightedSum - trainingData.Output[exampleIndex]
			for featureIndex := 0; featureIndex < len(trainingData.Input[exampleIndex]); featureIndex++ {
				model.Weights[featureIndex] = model.Weights[featureIndex] - model.LearningRate*loss*trainingData.Input[exampleIndex][featureIndex]
			}
		}
	}
	model.Status = "trained"
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model)
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	modelId := r.URL.Path[len("/test/"):]
	newData := &DataSet{
		Name:   "demo-dataset",
		Input:  [][]float64{{0.4, 0.3}, {0.7, 0.2}},
		Output: []float64{1.5, 2.3},
	}
	mu.Lock()
	model, ok := models[modelId]
	if !ok {
		mu.Unlock()
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}
	loss := 0.0
	for exampleIndex := 0; exampleIndex < len(newData.Input); exampleIndex++ {
		weightedSum := 0.0
		for featureIndex := 0; featureIndex < len(newData.Input[exampleIndex]); featureIndex++ {
			weightedSum += model.Weights[featureIndex] * newData.Input[exampleIndex][featureIndex]
		}
		loss += weightedSum - newData.Output[exampleIndex]
	}
	model.Loss = loss
	model.Status = "tested"
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
	http.HandleFunc("/train/", handleTrain)
	http.HandleFunc("/test/", handleTest)
	http.HandleFunc("/status/", handleStatus)
	http.HandleFunc("/models", handleModels)
	http.HandleFunc("/health", handleHealth)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
