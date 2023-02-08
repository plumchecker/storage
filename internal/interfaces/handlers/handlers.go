package handlers

import (
	"encoding/json"
	"net/http"

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

		counter, err := app.AddLeaks(newLeaks.Leaks)
		if err != nil {
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(counter)
	})
}

func getLeaks(app storage.Controller) http.Handler {
	type Keyword struct {
		Key   string
		Value string
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error getting leaks"
		newKeyword := new(Keyword)
		err := json.NewDecoder(r.Body).Decode(&newKeyword)

		leaks, err := app.FindLeaksByKeyword(newKeyword.Key, newKeyword.Value)
		if err != nil {
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(leaks)
	})
}

func Make(r *mux.Router, app storage.Controller) {
	apiURI := "/api"
	serviceRouter := r.PathPrefix(apiURI).Subrouter()
	serviceRouter.Handle("/leaks", addLeaks(app)).Methods("POST")
	serviceRouter.Handle("/search", getLeaks(app)).Methods("POST")
}
