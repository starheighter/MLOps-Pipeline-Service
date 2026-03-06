package handlersservices

type DataSet struct {
	Code   string      `json:"code"`
	Name   string      `json:"name"`
	Input  [][]float64 `json:"input"`
	Output []float64   `json:"output"`
}

func Trainset(code string) *DataSet {
	trainData := &DataSet{
		Code:   "dy",
		Name:   "Dummy",
		Input:  [][]float64{{0.0, 0.0}},
		Output: []float64{0.0},
	}
	switch code {
	case "hp":
		trainData = &DataSet{
			Code:   "hp",
			Name:   "HousingPrices",
			Input:  [][]float64{{0.2, 0.4}, {0.1, 0.8}, {0.3, 0.3}, {0.5, 0.2}},
			Output: []float64{1.0, 1.1, 1.2, 1.7},
		}
	case "rd":
		trainData = &DataSet{
			Code:   "rd",
			Name:   "RetailDemand",
			Input:  [][]float64{{0.2, 0.4}, {0.1, 0.8}, {0.3, 0.3}, {0.5, 0.2}},
			Output: []float64{0.0, -0.6, 0.3, 0.8},
		}
	case "tf":
		trainData = &DataSet{
			Code:   "tf",
			Name:   "TrafficFlow",
			Input:  [][]float64{{0.2, 0.4}, {0.1, 0.8}, {0.3, 0.3}, {0.5, 0.2}},
			Output: []float64{0.2, 0.7, 0.0, -0.3},
		}
	}
	return trainData
}

func Testset(code string) *DataSet {
	testData := &DataSet{
		Code:   "dy",
		Name:   "Dummy",
		Input:  [][]float64{{0.0, 0.0}},
		Output: []float64{0.0},
	}
	switch code {
	case "hp":
		testData = &DataSet{
			Code:   "hp",
			Name:   "HousingPrices",
			Input:  [][]float64{{0.4, 0.3}, {0.2, 0.7}},
			Output: []float64{1.5, 1.3},
		}
	case "rd":
		testData = &DataSet{
			Code:   "rd",
			Name:   "RetailDemand",
			Input:  [][]float64{{0.4, 0.3}, {0.2, 0.7}},
			Output: []float64{0.5, -0.3},
		}
	case "tf":
		testData = &DataSet{
			Code:   "tf",
			Name:   "TrafficFlow",
			Input:  [][]float64{{0.4, 0.3}, {0.2, 0.7}},
			Output: []float64{-0.1, 0.5},
		}
	}
	return testData
}
