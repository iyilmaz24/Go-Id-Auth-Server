package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/iyilmaz24/Go-Id-Auth-Server/internal/database/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "Welcome to the Golang Server Catch-All")
	fmt.Fprintln(w, "Use Correct Routes and Methods.")
}

func (app *application) getSurvey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		ID string `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.ID == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	cleanedId := strings.ToUpper(input.ID)

	survey, err := app.surveys.Get(cleanedId)
	if err != nil {
		if err == models.ErrNoRecord {
			app.clientError(w, http.StatusNotFound)
			return
		}
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(survey)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) getSurveyDbHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats, err := app.surveys.CheckHealth()
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		app.errorLog.Printf("Could not encode health check response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

