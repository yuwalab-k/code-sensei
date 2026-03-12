package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExecuteRequest struct {
	Code string `json:"code"`
}

type ExecuteResponse struct {
	Output string `json:"output"`
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request at /execute")
	var req ExecuteRequest
	json.NewDecoder(r.Body).Decode(&req)
	res := ExecuteResponse{Output: "not implemented yet"}
	json.NewEncoder(w).Encode(res)
}

func main() {
	fmt.Println("Server starting on 0.0.0.0:8081...")
	http.HandleFunc("/execute", executeHandler)
	err := http.ListenAndServe("0.0.0.0:8081", nil)
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}