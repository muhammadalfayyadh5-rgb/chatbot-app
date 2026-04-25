package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	
	"google.golang.org/genai"
)

type ChatRequest struct {
	Message string `json:"message"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	// PENTING: Izinkan frontend temanmu mengakses API ini
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	// Ambil input dari frontend
	var req ChatRequest
	json.NewDecoder(r.Body).Decode(&req)

	apiKey := os.Getenv("GEMINI_API_KEY")
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})

	// Kirim pesan dari user ke Gemini
	resp, err := client.Models.GenerateContent(ctx, "modeld/gemini2.5-flash-lite", genai.Text(req.Message), nil)
	
	botMessage := "Gagal mendapatkan respon."
	if err == nil && len(resp.Candidates) > 0 {
		botMessage = resp.Candidates[0].Content.Parts[0].Text
	}

	json.NewEncoder(w).Encode(map[string]string{"reply": botMessage})
}

func main() {
	http.HandleFunc("/chat", chatHandler)
	http.ListenAndServe(":8080", nil)
}