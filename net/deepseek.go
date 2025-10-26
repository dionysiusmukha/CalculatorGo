package net

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const proxyURL = "https://deproxy.kchugalinskiy.ru/deeproxy/api/completions"
const username = "41-2"
const password = "U0dMUjFs"

func SendToDeepSeek(query string) (string, error) {
	systemPrompt := "You are a knowledgeable assistant that provides accurate, reliable, and meaningful responses to user queries with high confidence."
	messages := []map[string]string{
		{"role": "system", "content": systemPrompt},
		{"role": "user", "content": query},
	}
	body := map[string]interface{}{
		"model":    "deepseek-chat",
		"messages": messages,
		"stream":   false,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", proxyURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	choice := choices[0].(map[string]interface{})
	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("no message in choice")
	}
	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("no content in message")
	}
	return content, nil
}



type Extracted struct {
	Action      string `json:"action"`      
	URL         string `json:"url,omitempty"`
	Instruction string `json:"instruction,omitempty"`
}


func sendMessages(messages []map[string]string) (string, error) {
	body := map[string]interface{}{
		"model":    "deepseek-chat",
		"messages": messages,
		"stream":   false,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", proxyURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	choice := choices[0].(map[string]interface{})
	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("no message in choice")
	}
	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("no content in message")
	}
	return content, nil
}


func ExtractFreeForm(userText string) (*Extracted, error) {
	system := `Ты — парсер команд. Преобразуй свободный текст пользователя в СТРОГО валидный JSON
c полями: {"action": "curl|calc|chat", "url": string|null, "instruction": string|null}.
- "curl": когда нужно открыть веб-страницу или "открой сайт/страницу/перейди по ссылке".
- "calc": когда это математическое выражение/операция.
- "chat": всё остальное (опрос, вопрос, просьба без URL).
ЕСЛИ присутствует URL (http/https), помести его в "url".
"instruction" — что сделать ПОСЛЕ получения содержимого (например, "сделай краткую сводку").
Выводи ТОЛЬКО JSON без комментариев и текста вокруг.`
	user := userText

	out, err := sendMessages([]map[string]string{
		{"role": "system", "content": system},
		{"role": "user", "content": user},
	})
	if err != nil {
		return nil, err
	}
	var e Extracted
	
	start := bytes.IndexByte([]byte(out), '{')
	end := bytes.LastIndexByte([]byte(out), '}')
	if start >= 0 && end > start {
		out = string([]byte(out)[start : end+1])
	}
	if err := json.Unmarshal([]byte(out), &e); err != nil {
		return nil, fmt.Errorf("parse extractor JSON failed: %v; raw: %s", err, out)
	}
	return &e, nil
}


func AskWithContext(instruction, context string) (string, error) {
	const max = 100 * 1024
	if len(context) > max {
		context = context[:max] + "\n\n[контент обрезан]"
	}
	system := "Ты — помощник, который выполняет инструкцию пользователя над переданным контекстом страницы. Пиши кратко и по делу."
	user := fmt.Sprintf("Инструкция: %s\n\n=== КОНТЕКСТ СТРАНИЦЫ НАЧАЛО ===\n%s\n=== КОНТЕКСТ СТРАНИЦЫ КОНЕЦ ===", instruction, context)
	return sendMessages([]map[string]string{
		{"role": "system", "content": system},
		{"role": "user", "content": user},
	})
}

