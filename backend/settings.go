package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const settingsPath = "/app/db/settings.json"

type Settings struct {
	BirthYear int `json:"birth_year"`
}

func loadSettings() Settings {
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return Settings{}
	}
	var s Settings
	json.Unmarshal(data, &s)
	return s
}

func saveSettings(s Settings) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(settingsPath, data, 0644)
}

func settingsGetHandler(w http.ResponseWriter, r *http.Request) {
	s := loadSettings()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func settingsUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var s Settings
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := saveSettings(s); err != nil {
		fmt.Printf("settings save error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}
