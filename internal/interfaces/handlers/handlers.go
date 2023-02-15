package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/D3vR4pt0rs/logger"
	"github.com/gorilla/mux"
	"github.com/plumchecker/storage/internal/entities"
	"github.com/plumchecker/storage/internal/usecases/storage"
)

func addLeaks(app storage.Controller) http.Handler {
	type Leaks struct {
		Leaks []entities.Leak
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error add leaks"
		newLeaks := new(Leaks)
		err := json.NewDecoder(r.Body).Decode(&newLeaks)

		logger.Info.Printf("Get request for add %d leaks", len(newLeaks.Leaks))

		counter, err := app.AddLeaks(newLeaks.Leaks)
		if err != nil {
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := make(map[string]interface{})
		response["counter"] = counter
		json.NewEncoder(w).Encode(response)
	})
}

func getLeaks(app storage.Controller) http.Handler {
	type Keyword struct {
		Key   string `json:"key"`
		Value string `json:"value,omitempty"`
		Token string `json:"token,omitempty"`
	}
	type Response struct {
		leaks []entities.Leak
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error getting leaks"

		newKeyword := new(Keyword)
		err := json.NewDecoder(r.Body).Decode(&newKeyword)
		logger.Info.Printf("Get request for getting values by %s with value %s", newKeyword.Key, newKeyword.Value)

		leaks, token, err := app.FindLeaksByKeyword(newKeyword.Key, newKeyword.Value, newKeyword.Token)
		if err != nil {
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := make(map[string]interface{})
		response["leaks"] = leaks
		response["token"] = token
		json.NewEncoder(w).Encode(response)
	})
}

func Make(r *mux.Router, app storage.Controller) {
	apiURI := "/api"
	serviceRouter := r.PathPrefix(apiURI).Subrouter()
	serviceRouter.Handle("/leaks", addLeaks(app)).Methods("POST")
	serviceRouter.Handle("/search", getLeaks(app)).Methods("POST")
}
