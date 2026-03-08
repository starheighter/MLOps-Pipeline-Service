package handlersservices

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Model struct {
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"createdAt"`
	NumberTrainLoops int       `json:"numberTrainLoops"`
	LearningRate     float64   `json:"learningRate"`
	Loss             float64   `json:"loss"`
	Weights          []float64 `json:"weights"`
	TrainInputFile   string    `json:"trainInputFile"`
	TrainOutputFile  string    `json:"trainOutputFile"`
	TestInputFile    string    `json:"testInputFile"`
	TestOutputFile   string    `json:"testOutputFile"`
}

var (
	models = make(map[string]*Model)
	mu     sync.Mutex
)

func CreateModel(modelName string, numberTrainLoops int, learningRate float64) *Model {
	model := &Model{
		Name:             modelName,
		Status:           "created",
		CreatedAt:        time.Now(),
		NumberTrainLoops: numberTrainLoops,
		LearningRate:     learningRate,
		Loss:             0.0,
		Weights:          []float64{},
		TrainInputFile:   "",
		TrainOutputFile:  "",
		TestInputFile:    "",
		TestOutputFile:   "",
	}
	models[model.Name] = model
	return model
}

func HandleTrain(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/train.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func HandleTest(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/test.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func HandleTraining(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(10 << 20)
	modelName := r.FormValue("model_name")
	numberTrainLoops, err := strconv.Atoi(r.FormValue("number_train_loops"))
	if err != nil {
		mu.Unlock()
		http.Error(w, "Number Train Loops must be a number", http.StatusBadRequest)
		return
	}
	learningRate, err := strconv.ParseFloat(r.FormValue("learning_rate"), 64)
	if err != nil {
		mu.Unlock()
		http.Error(w, "Learning Rate must be a float", http.StatusBadRequest)
		return
	}
	mu.Lock()
	model := CreateModel(modelName, numberTrainLoops, learningRate)
	trainData := ExtractDataSet(w, r)
	if trainData == nil {
		mu.Unlock()
		return
	}
	model.TrainInputFile = trainData.InputFile
	model.TrainOutputFile = trainData.OutputFile
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
	modelName := r.FormValue("model_name")
	mu.Lock()
	model, ok := models[modelName]
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
	model.TestInputFile = testData.InputFile
	model.TestOutputFile = testData.OutputFile
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
