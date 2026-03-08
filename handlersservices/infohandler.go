package handlersservices

import (
	"encoding/json"
	"net/http"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	var list []*Model
	for _, m := range models {
		list = append(list, m)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func HandleModel(w http.ResponseWriter, r *http.Request) {
	modelName := r.URL.Path[len("/model/"):]
	mu.Lock()
	model, ok := models[modelName]
	mu.Unlock()
	if !ok {
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model)
}
