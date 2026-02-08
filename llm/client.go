package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client interface {
	ReviewCode(code string) (string, error)
}

type OpenAIClient struct {
	APIKey string
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{APIKey: apiKey}
}

func (c *OpenAIClient) ReviewCode(code string) (string, error) {
	if c.APIKey == "" {
		return "Skipped LLM review: OpenAI API key not set.", nil
	}

	prompt := fmt.Sprintf("Review the following Go code for bugs, security issues, and performance improvements. Be concise.\n\n%s", code)

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]string{
			{"role": "system", "content": "You are an expert Go code reviewer."},
			{"role": "user", "content": prompt},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error: %s", string(body))
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]any)
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("invalid response from OpenAI")
	}

	message := choices[0].(map[string]any)["message"].(map[string]interface{})
	return message["content"].(string), nil
}

type OllamaClient struct {
	Model string
}

func NewOllamaClient(model string) *OllamaClient {
	return &OllamaClient{Model: model}
}

func (c *OllamaClient) ReviewCode(code string) (string, error) {
	prompt := fmt.Sprintf("Review the following Go code for bugs, security issues, and performance improvements. Be concise.\n\n%s", code)

	requestBody, err := json.Marshal(map[string]any{
		"model":  c.Model,
		"prompt": prompt,
		"stream": false,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "http://localhost:11434/api/generate", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API error: %s", string(body))
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	response, ok := result["response"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response from Ollama")
	}

	return response, nil
}
