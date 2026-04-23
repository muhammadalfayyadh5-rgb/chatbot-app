package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Reply string `json:"reply"`
}

// 🔥 Fungsi panggil Gemini API
func callGeminiAPI(message string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", apiKey)

	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": message},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	resBody, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(resBody, &result)

	// 🔍 Ambil teks dari response Gemini
	candidates := result["candidates"].([]interface{})
	content := candidates[0].(map[string]interface{})["content"].(map[string]interface{})
	parts := content["parts"].([]interface{})
	text := parts[0].(map[string]interface{})["text"].(string)

	return text, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Request tidak valid", http.StatusBadRequest)
		return
	}

	// 🔥 Panggil Gemini
	reply, err := callGeminiAPI(req.Message)
	if err != nil {
		http.Error(w, "Gagal memanggil AI", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	res := Response{Reply: reply}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

	log.Println("User:", req.Message)
	log.Println("AI:", reply)
}

func main() {
	http.HandleFunc("/chat", handler)

	log.Println("🚀 Server AI berjalan di port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}