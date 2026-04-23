package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Reply string `json:"reply"`
}

func main() {
	req := Request{Message: "halo"}
	data, _ := json.Marshal(req)

	resp, err := http.Post("http://192.168.1.5:8080/chat", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var res Response
	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println("Bot:", res.Reply)
}