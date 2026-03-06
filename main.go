package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/starheighter/MLOpsPipelineservice/handlersservices"
)

func main() {
	if os.Getenv("SEED_JOBS") == "true" {
		for i := 0; i < 5; i++ {
			handlersservices.CreateModel()
		}
	}
	http.HandleFunc("/", handlersservices.HandleHome)
	http.HandleFunc("/train/", handlersservices.HandleTrain)
	http.HandleFunc("/test/", handlersservices.HandleTest)
	http.HandleFunc("/model/", handlersservices.HandleModel)
	http.HandleFunc("/models", handlersservices.HandleModels)
	http.HandleFunc("/health", handlersservices.HandleHealth)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
