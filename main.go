package main

import (
	"fmt"
	"net/http"

	"github.com/starheighter/MLOpsPipelineservice/handlersservices"
)

func main() {
	http.HandleFunc("/", handlersservices.HandleHome)
	http.HandleFunc("/train", handlersservices.HandleTrain)
	http.HandleFunc("/test", handlersservices.HandleTest)
	http.HandleFunc("/training/", handlersservices.HandleTraining)
	http.HandleFunc("/testing/", handlersservices.HandleTesting)
	http.HandleFunc("/model/", handlersservices.HandleModel)
	http.HandleFunc("/health", handlersservices.HandleHealth)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
