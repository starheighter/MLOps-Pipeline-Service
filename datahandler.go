package main

func Trainset(index string) *DataSet {
	trainData := &DataSet{
		Name:   "dummytestset",
		Input:  [][]float64{{0.0, 0.0}},
		Output: []float64{0.0},
	}
	switch index {
	case "1":
		trainData = &DataSet{
			Name:   "trainset1",
			Input:  [][]float64{{0.2, 0.4}, {0.1, 0.8}, {0.3, 0.3}, {0.5, 0.2}},
			Output: []float64{1.0, 1.1, 1.2, 1.7},
		}
	case "2":
		trainData = &DataSet{
			Name:   "trainset2",
			Input:  [][]float64{{0.2, 0.4}, {0.1, 0.8}, {0.3, 0.3}, {0.5, 0.2}},
			Output: []float64{0.0, -0.6, 0.3, 0.8},
		}
	case "3":
		trainData = &DataSet{
			Name:   "trainset3",
			Input:  [][]float64{{0.2, 0.4}, {0.1, 0.8}, {0.3, 0.3}, {0.5, 0.2}},
			Output: []float64{0.2, 0.7, 0.0, -0.3},
		}
	}
	return trainData
}

func Testset(index string) *DataSet {
	testData := &DataSet{
		Name:   "dummytestset",
		Input:  [][]float64{{0.0, 0.0}},
		Output: []float64{0.0},
	}
	switch index {
	case "1":
		testData = &DataSet{
			Name:   "testset1",
			Input:  [][]float64{{0.4, 0.3}, {0.2, 0.7}},
			Output: []float64{1.5, 1.3},
		}
	case "2":
		testData = &DataSet{
			Name:   "testset2",
			Input:  [][]float64{{0.4, 0.3}, {0.2, 0.7}},
			Output: []float64{0.5, -0.3},
		}
	case "3":
		testData = &DataSet{
			Name:   "testset3",
			Input:  [][]float64{{0.4, 0.3}, {0.2, 0.7}},
			Output: []float64{-0.1, 0.5},
		}
	}
	return testData
}
