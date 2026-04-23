package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Reply string `json:"reply"`
}

func chatbot(msg string) string {
	msg = strings.ToLower(msg)

	if strings.Contains(msg, "halo") {
		return "Halo juga!"
	}
	if strings.Contains(msg, "siapa kamu") {
		return "Saya chatbot Golang 🤖"
	}
	return "Saya tidak mengerti"
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req Request
	json.NewDecoder(r.Body).Decode(&req)

	reply := chatbot(req.Message)

	res := Response{Reply: reply}
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/chat", handler)
	http.ListenAndServe(":8080", nil)
}