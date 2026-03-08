package handlersservices

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
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
	InputFile  string      `json:"inputHeader"`
	OutputFile string      `json:"outputHeader"`
	Input      [][]float64 `json:"input"`
	Output     []float64   `json:"output"`
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
		Weights:          []float64{0.0, 0.0},
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
	inputFile, inputHeader, err := r.FormFile("input_file")
	if err != nil {
		mu.Unlock()
		http.Error(w, "Input file missing", http.StatusBadRequest)
		return
	}
	defer inputFile.Close()
	outputFile, outputHeader, err := r.FormFile("output_file")
	if err != nil {
		mu.Unlock()
		http.Error(w, "Output file missing", http.StatusBadRequest)
		return
	}
	defer outputFile.Close()
	var inputData [][]float64
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Fields(line)
		var row []float64
		for _, v := range values {
			num, err := strconv.ParseFloat(v, 64)
			if err != nil {
				mu.Unlock()
				http.Error(w, "Input data has to contain only float numbers", http.StatusBadRequest)
				return
			}
			row = append(row, num)
		}
		inputData = append(inputData, row)
	}
	var outputData []float64
	scanner = bufio.NewScanner(outputFile)
	for scanner.Scan() {
		line := scanner.Text()
		num, err := strconv.ParseFloat(line, 64)
		if err != nil {
			mu.Unlock()
			http.Error(w, "Output data has to contain only float numbers", http.StatusBadRequest)
			return
		}
		outputData = append(outputData, num)
	}
	trainData := &DataSet{
		InputFile:  inputHeader.Filename,
		OutputFile: outputHeader.Filename,
		Input:      inputData,
		Output:     outputData,
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
	inputFile, inputHeader, err := r.FormFile("input_file")
	if err != nil {
		mu.Unlock()
		http.Error(w, "Input file missing", http.StatusBadRequest)
		return
	}
	defer inputFile.Close()
	outputFile, outputHeader, err := r.FormFile("output_file")
	if err != nil {
		mu.Unlock()
		http.Error(w, "Output file missing", http.StatusBadRequest)
		return
	}
	defer outputFile.Close()
	var inputData [][]float64
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Fields(line)
		var row []float64
		for _, v := range values {
			num, err := strconv.ParseFloat(v, 64)
			if err != nil {
				mu.Unlock()
				http.Error(w, "Input data has to contain only float numbers", http.StatusBadRequest)
				return
			}
			row = append(row, num)
		}
		inputData = append(inputData, row)
	}
	var outputData []float64
	scanner = bufio.NewScanner(outputFile)
	for scanner.Scan() {
		line := scanner.Text()
		num, err := strconv.ParseFloat(line, 64)
		if err != nil {
			mu.Unlock()
			http.Error(w, "Output data has to contain only float numbers", http.StatusBadRequest)
			return
		}
		outputData = append(outputData, num)
	}
	testData := &DataSet{
		InputFile:  inputHeader.Filename,
		OutputFile: outputHeader.Filename,
		Input:      inputData,
		Output:     outputData,
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
	http.Redirect(w, r, "/", http.StatusFound)
}
