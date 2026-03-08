package handlersservices

import (
	"bufio"
	"net/http"
	"strconv"
	"strings"
)

type DataSet struct {
	InputFile  string      `json:"inputFile"`
	OutputFile string      `json:"outputFile"`
	InputData  [][]float64 `json:"inputData"`
	OutputData []float64   `json:"outputData"`
}

func ExtractDataSet(w http.ResponseWriter, r *http.Request) *DataSet {
	inputFile, inputHeader, err := r.FormFile("input_file")
	if err != nil {
		http.Error(w, "Input file missing", http.StatusBadRequest)
		return nil
	}
	defer inputFile.Close()
	outputFile, outputHeader, err := r.FormFile("output_file")
	if err != nil {
		http.Error(w, "Output file missing", http.StatusBadRequest)
		return nil
	}
	defer outputFile.Close()
	var inputData [][]float64
	scanner := bufio.NewScanner(inputFile)
	counterInputs := 0
	for scanner.Scan() {
		counterInputs++
		line := scanner.Text()
		values := strings.Fields(line)
		var row []float64
		for _, v := range values {
			num, err := strconv.ParseFloat(v, 64)
			if err != nil {
				http.Error(w, "Input data has to contain only float numbers", http.StatusBadRequest)
				return nil
			}
			row = append(row, num)
		}
		inputData = append(inputData, row)
	}
	var outputData []float64
	scanner = bufio.NewScanner(outputFile)
	counterOutputs := 0
	for scanner.Scan() {
		counterOutputs++
		line := scanner.Text()
		num, err := strconv.ParseFloat(line, 64)
		if err != nil {
			http.Error(w, "Output data has to contain only float numbers", http.StatusBadRequest)
			return nil
		}
		outputData = append(outputData, num)
	}
	if counterInputs == 0 || counterInputs != counterOutputs {
		http.Error(w, "Illegal number of input / output samples", http.StatusBadRequest)
		return nil
	}
	dataSet := &DataSet{
		InputFile:  inputHeader.Filename,
		OutputFile: outputHeader.Filename,
		InputData:  inputData,
		OutputData: outputData,
	}
	return dataSet
}
