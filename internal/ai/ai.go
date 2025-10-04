package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/justynroberts/finpup/internal/config"
)

type Client struct {
	config *config.AIConfig
}

func New(cfg *config.AIConfig) *Client {
	return &Client{config: cfg}
}

func (c *Client) GenerateText(prompt string, context string) (string, error) {
	if !c.config.Enabled {
		return "", fmt.Errorf("AI is not enabled in config")
	}

	switch c.config.Provider {
	case "ollama":
		return c.generateOllama(prompt, context)
	case "openai":
		return c.generateOpenAI(prompt, context)
	case "openrouter":
		return c.generateOpenRouter(prompt, context)
	default:
		return "", fmt.Errorf("unsupported AI provider: %s", c.config.Provider)
	}
}

func (c *Client) generateOllama(prompt string, context string) (string, error) {
	fullPrompt := prompt
	if context != "" {
		fullPrompt = fmt.Sprintf("Context:\n%s\n\nTask: %s", context, prompt)
	}

	reqBody := map[string]interface{}{
		"model":  c.config.Model,
		"prompt": fullPrompt,
		"stream": false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := strings.TrimSuffix(c.config.BaseURL, "/") + "/api/generate"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API error (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return strings.TrimSpace(result.Response), nil
}

func (c *Client) generateOpenAI(prompt string, context string) (string, error) {
	messages := []map[string]string{}

	if context != "" {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": "You are a helpful code assistant. Provide concise responses.",
		})
		messages = append(messages, map[string]string{
			"role":    "user",
			"content": fmt.Sprintf("Context:\n%s\n\nTask: %s", context, prompt),
		})
	} else {
		messages = append(messages, map[string]string{
			"role":    "user",
			"content": prompt,
		})
	}

	reqBody := map[string]interface{}{
		"model":    c.config.Model,
		"messages": messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := strings.TrimSuffix(c.config.BaseURL, "/") + "/v1/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

func (c *Client) generateOpenRouter(prompt string, context string) (string, error) {
	messages := []map[string]string{}

	if context != "" {
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": "You are a helpful code assistant. Provide concise responses.",
		})
		messages = append(messages, map[string]string{
			"role":    "user",
			"content": fmt.Sprintf("Context:\n%s\n\nTask: %s", context, prompt),
		})
	} else {
		messages = append(messages, map[string]string{
			"role":    "user",
			"content": prompt,
		})
	}

	reqBody := map[string]interface{}{
		"model":    c.config.Model,
		"messages": messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := strings.TrimSuffix(c.config.BaseURL, "/") + "/api/v1/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("HTTP-Referer", "https://github.com/justynroberts/finpup")
	req.Header.Set("X-Title", "finpup Editor")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenRouter API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenRouter API error (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}
