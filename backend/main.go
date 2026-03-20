package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server starting on 0.0.0.0:8081...")
	http.HandleFunc("/execute", cors(executeHandler))
	http.HandleFunc("/settings", cors(settingsGetHandler))
	http.HandleFunc("/settings/update", cors(settingsUpdateHandler))
	http.HandleFunc("/ai/hint", cors(aiHintHandler))
	http.HandleFunc("/ai/generate", cors(aiGenerateHandler))
	http.HandleFunc("/ai/review", cors(aiReviewHandler))
	http.HandleFunc("/ai/blank", cors(aiBlankHandler))
	if err := http.ListenAndServe("0.0.0.0:8081", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
