package handlersservices

import (
	"encoding/json"
	"net/http"
	"sync"
)

var (
	models = make(map[string]*Model)
	mu     sync.Mutex
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/models", http.StatusFound)
}

func HandleModel(w http.ResponseWriter, r *http.Request) {
	modelId := r.URL.Path[len("/model/"):]
	mu.Lock()
	model, ok := models[modelId]
	mu.Unlock()
	if !ok {
		http.Error(w, "Model not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model)
}

func HandleModels(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	var list []*Model
	for _, m := range models {
		list = append(list, m)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
