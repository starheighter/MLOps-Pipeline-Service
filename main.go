package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Model struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
	NumberTrainLoops int       `json:"numberTrainLoops"`
	LearningRate     float64   `json:"learningRate"`
	Loss             float64   `json:"loss"`
	Weights          []float64 `json:"weights"`
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
	modelId := fmt.Sprintf("%d", 10000+rand.Intn(90000))
	model := &Model{
		ID:               modelId,
		Name:             "demo-model",
		Status:           "created",
		CreatedAt:        time.Now(),
		NumberTrainLoops: 150,
		LearningRate:     0.05,
		Loss:             0.0,
		Weights:          []float64{0.0, 0.0},
	}
	mu.Lock()
	models[modelId] = model
	mu.Unlock()
	return model
}

func handleTrain(w http.ResponseWriter, r *http.Request) {
	parameters := strings.Split(r.URL.Path[len("/train/"):], "/")
	modelId := parameters[0]
	datasetIndex := parameters[1]
	trainData := Trainset(datasetIndex)
	mu.Lock()
	model, ok := models[modelId]
	if !ok {
		mu.Unlock()
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}
	for trainLoop := 0; trainLoop < model.NumberTrainLoops; trainLoop++ {
		for exampleIndex := 0; exampleIndex < len(trainData.Input); exampleIndex++ {
			weightedSum := 0.0
			loss := 0.0
			for featureIndex := 0; featureIndex < len(trainData.Input[exampleIndex]); featureIndex++ {
				weightedSum += model.Weights[featureIndex] * trainData.Input[exampleIndex][featureIndex]
			}
			loss += weightedSum - trainData.Output[exampleIndex]
			for featureIndex := 0; featureIndex < len(trainData.Input[exampleIndex]); featureIndex++ {
				model.Weights[featureIndex] = model.Weights[featureIndex] - model.LearningRate*loss*trainData.Input[exampleIndex][featureIndex]
			}
		}
	}
	model.Status = "trained"
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model)
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	parameters := strings.Split(r.URL.Path[len("/test/"):], "/")
	modelId := parameters[0]
	datasetIndex := parameters[1]
	testData := Testset(datasetIndex)
	mu.Lock()
	model, ok := models[modelId]
	if !ok {
		mu.Unlock()
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}
	loss := 0.0
	for exampleIndex := 0; exampleIndex < len(testData.Input); exampleIndex++ {
		weightedSum := 0.0
		for featureIndex := 0; featureIndex < len(testData.Input[exampleIndex]); featureIndex++ {
			weightedSum += model.Weights[featureIndex] * testData.Input[exampleIndex][featureIndex]
		}
		loss += weightedSum - testData.Output[exampleIndex]
	}
	model.Loss = loss
	model.Status = "tested"
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model)
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
	http.HandleFunc("/model/", handleModel)
	http.HandleFunc("/models", handleModels)
	http.HandleFunc("/health", handleHealth)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
