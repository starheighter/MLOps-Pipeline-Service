package handlersservices

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
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

var (
	models = make(map[string]*Model)
	mu     sync.Mutex
)

func CreateModel() *Model {
	modelId := fmt.Sprintf("%d", 10000+rand.Intn(90000))
	model := &Model{
		ID:               modelId,
		Name:             "demo-model",
		Status:           "created",
		CreatedAt:        time.Now(),
		NumberTrainLoops: 0,
		LearningRate:     0.0,
		Loss:             0.0,
		Weights:          []float64{},
	}
	mu.Lock()
	models[modelId] = model
	mu.Unlock()
	return model
}

func HandleTrain(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/train.html")
}

func HandleTest(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/test.html")
}

func HandleTraining(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(10 << 20)
	modelId := r.FormValue("model_id")
	mu.Lock()
	model, ok := models[modelId]
	if !ok {
		mu.Unlock()
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}
	numberTrainLoops, err := strconv.Atoi(r.FormValue("number_train_loops"))
	if err != nil {
		mu.Unlock()
		http.Error(w, "Number Train Loops must be a number", http.StatusBadRequest)
		return
	}
	model.NumberTrainLoops = numberTrainLoops
	learningRate, err := strconv.ParseFloat(r.FormValue("learning_rate"), 64)
	if err != nil {
		mu.Unlock()
		http.Error(w, "Learning Rate must be a float", http.StatusBadRequest)
		return
	}
	model.LearningRate = learningRate
	trainData := ExtractDataSet(w, r)
	if trainData == nil {
		mu.Unlock()
		return
	}
	model.Weights = make([]float64, len(trainData.InputData[0]))
	for trainLoop := 0; trainLoop < model.NumberTrainLoops; trainLoop++ {
		for sampleIndex := 0; sampleIndex < len(trainData.InputData); sampleIndex++ {
			weightedSum := 0.0
			loss := 0.0
			for featureIndex := 0; featureIndex < len(trainData.InputData[sampleIndex]); featureIndex++ {
				weightedSum += model.Weights[featureIndex] * trainData.InputData[sampleIndex][featureIndex]
			}
			loss += weightedSum - trainData.OutputData[sampleIndex]
			for featureIndex := 0; featureIndex < len(trainData.InputData[sampleIndex]); featureIndex++ {
				model.Weights[featureIndex] = model.Weights[featureIndex] - model.LearningRate*loss*trainData.InputData[sampleIndex][featureIndex]
			}
		}
	}
	model.Status = "trained"
	mu.Unlock()
	http.Redirect(w, r, "/", http.StatusFound)
}

func HandleTesting(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(10 << 20)
	modelId := r.FormValue("model_id")
	mu.Lock()
	model, ok := models[modelId]
	if !ok {
		mu.Unlock()
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}
	testData := ExtractDataSet(w, r)
	if testData == nil {
		mu.Unlock()
		return
	}
	loss := 0.0
	for sampleIndex := 0; sampleIndex < len(testData.InputData); sampleIndex++ {
		weightedSum := 0.0
		for featureIndex := 0; featureIndex < len(testData.InputData[sampleIndex]); featureIndex++ {
			weightedSum += model.Weights[featureIndex] * testData.InputData[sampleIndex][featureIndex]
		}
		loss += weightedSum - testData.OutputData[sampleIndex]
	}
	model.Loss = loss
	model.Status = "tested"
	mu.Unlock()
	http.Redirect(w, r, "/", http.StatusFound)
}
