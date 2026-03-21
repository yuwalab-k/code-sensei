package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ===== Request / Response types =====

type ExecuteRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type ExecuteResponse struct {
	Output string `json:"output"`
}

type AIHintRequest struct {
	Code       string `json:"code"`
	Question   string `json:"question"`
	Mission    string `json:"mission"`
	Language   string `json:"language"`
	ReviewMode bool   `json:"review_mode"`
	BirthYear  int    `json:"birth_year"`
}

type AIHintResponse struct {
	Hint string `json:"hint"`
}

type AIGenerateRequest struct {
	Description string `json:"description"`
	Language    string `json:"language"`
	ExecMode    string `json:"exec_mode"` // "browser" or "file"
	BirthYear   int    `json:"birth_year"`
}

type AIGenerateResponse struct {
	Code        string `json:"code"`
	Explanation string `json:"explanation"`
}

type AIReviewRequest struct {
	Code      string `json:"code"`
	Language  string `json:"language"`
	BirthYear int    `json:"birth_year"`
}

type AIReviewResponse struct {
	Questions string `json:"questions"`
}

type AIBlankRequest struct {
	Code      string `json:"code"`
	Language  string `json:"language"`
	BirthYear int    `json:"birth_year"`
}

type AIBlankResponse struct {
	BlankedCode string `json:"blanked_code"`
}

// ===== Handlers =====

func executeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request at /execute")
	var req ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ExecuteResponse{Output: "リクエストエラー"})
		return
	}

	if req.Language == "" {
		req.Language = "python"
	}

	body, _ := json.Marshal(req)
	resp, err := http.Post("http://runner:8002/run", "application/json", bytes.NewBuffer(body))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ExecuteResponse{Output: "実行サービスに接続できませんでした"})
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

func aiHintHandler(w http.ResponseWriter, r *http.Request) {
	var req AIHintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Language == "" {
		req.Language = "python"
	}

	userMsg := req.Question
	if userMsg == "" {
		userMsg = "このコードで困っています。ヒントをください。"
	}
	if req.Mission != "" {
		userMsg = fmt.Sprintf("ミッション「%s」に取り組んでいます。%s", req.Mission, userMsg)
	}
	if req.Code != "" {
		userMsg += fmt.Sprintf("\n\n今のコード:\n```%s\n%s\n```", req.Language, req.Code)
	}

	systemPrompt := buildHintPrompt(req.Language, req.ReviewMode, req.BirthYear)
	hint, err := callOllama(r.Context(), systemPrompt, userMsg)
	if err != nil {
		fmt.Printf("Ollama error: %v\n", err)
		hint = "今ちょっと考え中だよ。もう一度聞いてみてね！"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AIHintResponse{Hint: hint})
}

func aiGenerateHandler(w http.ResponseWriter, r *http.Request) {
	var req AIGenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Language == "" {
		req.Language = "python"
	}
	if req.ExecMode == "" {
		req.ExecMode = "browser"
	}

	systemPrompt := buildGeneratePrompt(req.Language, req.ExecMode, req.BirthYear)
	raw, err := callOllama(r.Context(), systemPrompt, req.Description)
	if err != nil {
		fmt.Printf("Ollama generate error: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AIGenerateResponse{Explanation: "うまく作れなかったよ。もう一度試してみてね！"})
		return
	}

	fmt.Printf("[generate] raw response (%d chars): %.300s\n", len(raw), raw)
	code, explanation := extractCodeBlock(raw)
	fmt.Printf("[generate] extracted code (%d chars), explanation (%d chars)\n", len(code), len(explanation))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AIGenerateResponse{Code: code, Explanation: explanation})
}

func aiReviewHandler(w http.ResponseWriter, r *http.Request) {
	var req AIReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Language == "" {
		req.Language = "python"
	}

	systemPrompt := buildReviewPrompt(req.Language, req.BirthYear)
	userMsg := fmt.Sprintf("このコードの復習問題を3つ作ってください。\n\n```%s\n%s\n```", req.Language, req.Code)

	questions, err := callOllama(r.Context(), systemPrompt, userMsg)
	if err != nil {
		fmt.Printf("Ollama review error: %v\n", err)
		questions = "問題を作るのに失敗したよ。もう一度試してみてね！"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AIReviewResponse{Questions: questions})
}

func aiBlankHandler(w http.ResponseWriter, r *http.Request) {
	var req AIBlankRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Language == "" {
		req.Language = "python"
	}

	systemPrompt := buildBlankPrompt(req.Language, req.BirthYear)
	userMsg := fmt.Sprintf("穴埋め問題を作ってください。\n\n```%s\n%s\n```", req.Language, req.Code)

	raw, err := callOllama(r.Context(), systemPrompt, userMsg)
	if err != nil {
		fmt.Printf("Ollama blank error: %v\n", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AIBlankResponse{BlankedCode: req.Code})
		return
	}

	blanked, _ := extractCodeBlock(raw)
	if blanked == "" {
		blanked = req.Code
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AIBlankResponse{BlankedCode: blanked})
}
