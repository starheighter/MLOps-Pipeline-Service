package handlersservices

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Model struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
	DatasetCode      string    `json:"datasetCode"`
	NumberTrainLoops int       `json:"numberTrainLoops"`
	LearningRate     float64   `json:"learningRate"`
	Loss             float64   `json:"loss"`
	Weights          []float64 `json:"weights"`
}

func CreateModel() *Model {
	modelId := fmt.Sprintf("%d", 10000+rand.Intn(90000))
	model := &Model{
		ID:               modelId,
		Name:             "demo-model",
		Status:           "created",
		CreatedAt:        time.Now(),
		DatasetCode:      "",
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

func HandleTrain(w http.ResponseWriter, r *http.Request) {
	parameters := strings.Split(r.URL.Path[len("/train/"):], "/")
	modelId := parameters[0]
	mu.Lock()
	model, ok := models[modelId]
	if !ok {
		mu.Unlock()
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}
	if model.DatasetCode == "" {
		if len(parameters) < 2 {
			mu.Unlock()
			http.Error(w, "Missing DatasetCode parameter", http.StatusBadRequest)
		}
		if parameters[1] != "dy" && parameters[1] != "hp" && parameters[1] != "rd" && parameters[1] != "tf" {
			mu.Unlock()
			http.Error(w, "DatasetCode not found", http.StatusNotFound)
			return
		}
		model.DatasetCode = parameters[1]
	}
	trainData := Trainset(model.DatasetCode)
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

func HandleTest(w http.ResponseWriter, r *http.Request) {
	modelId := r.URL.Path[len("/test/"):]
	mu.Lock()
	model, ok := models[modelId]
	if !ok {
		mu.Unlock()
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}
	testData := Testset(model.DatasetCode)
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
