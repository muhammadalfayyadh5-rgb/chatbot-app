package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"google.golang.org/genai"
)

type ChatRequest struct {
	Message string `json:"message"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	// 1. CORS Headers - Sangat penting agar Frontend (port 80) bisa akses Backend (port 8080)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 2. Decode Input dari Frontend
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Request body tidak valid", http.StatusBadRequest)
		return
	}

	// 3. Ambil API Key dari Environment Variable
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("Error: Variable GEMINI_API_KEY tidak ditemukan!")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"reply": "Server error: API Key belum dikonfigurasi."})
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Printf("Gagal inisialisasi client: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"reply": "Gagal terhubung ke layanan AI."})
		return
	}

	// 4. Generate Content (Penyebab Error 404 Sebelumnya)
	// Di SDK ini, cukup tulis nama modelnya saja tanpa prefix "models/"
	modelID := "gemini-2.5-flash-lite" // Pastikan model ini sudah tersedia di akun Anda

	resp, err := client.Models.GenerateContent(ctx, modelID, genai.Text(req.Message), nil)

	if err != nil {
		log.Printf("Error Gemini: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		// Memberikan info error yang lebih manusiawi di chat
		json.NewEncoder(w).Encode(map[string]string{
			"reply": fmt.Sprintf("Waduh, Google Gemini bilang: %v", err),
		})
		return
	}

	// 5. Parse Response - Ambil teks dari hasil AI
	if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		json.NewEncoder(w).Encode(map[string]string{"reply": "Gemini tidak memberikan respon."})
		return
	}

	var answer strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			answer.WriteString(part.Text)
		}
	}

	botMessage := answer.String()
	if botMessage == "" {
		botMessage = "Maaf, jawaban tidak dapat diproses."
	}

	// Kirim balik ke Frontend
	json.NewEncoder(w).Encode(map[string]string{"reply": botMessage})
}

func main() {
	port := ":8080"
	http.HandleFunc("/chat", chatHandler)
	log.Printf("Backend AI Chatbot berjalan di port %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server gagal berjalan: ", err)
	}
}
