package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func Generate() {
	requestData := struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}{
		Model:  "deepseek-r1:1.5b",
		Prompt: "天空为什么是蓝色?",
	}

	// 序列化JSON数据
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Fatal("JSON marshal error:", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", "http://localhost:11434/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Request creation error:", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Request failed:", err)
	}
	defer resp.Body.Close()

	log.Printf("Res content type: %s\n", resp.Header.Get("Content-Type"))
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var chunk GenerateChunk
		err = json.Unmarshal(scanner.Bytes(), &chunk)
		if err != nil {
			log.Printf("JSON unmarshal error: %v\n", err)
			continue
		}
		log.Printf("%-20sDone: %t\n", chunk.Response, chunk.Done)
		if chunk.Done {
			log.Printf("Done reason: %s\n", chunk.DoneReason)
		}
	}
}

func GetTags() {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	var res ModelResponse
	json.Unmarshal(body, &res)
	log.Printf("Models: %+v\n", res.Models)
}

func Chat() {
	// TODO: Implement chat functionality
}
