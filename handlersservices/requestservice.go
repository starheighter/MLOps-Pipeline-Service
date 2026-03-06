package handlersservices

type TrainRequest struct {
	ModelName        string  `json:"modelName"`
	DatasetCode      string  `json:"datasetCode"`
	NumberTrainLoops int     `json:"numberTrainLoops"`
	LearningRate     float64 `json:"learningRate"`
}

type TestRequest struct {
	ModelName string `json:"modelName"`
}
