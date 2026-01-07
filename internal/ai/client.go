package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	MistralBaseURL = "https://api.mistral.ai/v1"
)

// TODO: need to do some cleanup to remove hard-coded use of Mistral Agents

type AIClient struct {
	apiEndpoint string
	apiKey      string
	modelID     string // model id or agent id
	httpClient  *http.Client
}

func NewAIClient(apiKey, mode, model string) (*AIClient, error) {
	if mode == "mistral-agent" {
		log.Println("Using Mistral Agent mode with agent ID:", model)

		return &AIClient{
			apiEndpoint: fmt.Sprintf("%s/agents/completions", MistralBaseURL),
			apiKey:      apiKey,
			modelID:     model,
			httpClient: &http.Client{
				Timeout: 60 * time.Second,
			},
		}, nil
	}

	if mode == "mistral-api" {
		log.Println("Using Mistral API mode with model ID:", model)

		return &AIClient{
			apiEndpoint: fmt.Sprintf("%s/chat/completions", MistralBaseURL),
			apiKey:      apiKey,
			modelID:     model,
			httpClient: &http.Client{
				Timeout: 60 * time.Second,
			},
		}, nil
	}

	return nil, fmt.Errorf("unsupported AI client mode: %s", mode)
}

// Chat sends a chat completion request AI API configured
func (c *AIClient) Chat(messages []AIMessage) (string, error) {
	return c.chatWithAgent(messages)
}

// chatWithAgent sends a request to the agent completion endpoint
func (c *AIClient) chatWithAgent(messages []AIMessage) (string, error) {
	reqBody := AgentRequest{
		AgentID:  c.modelID,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, errResp.Error.Message)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// ChatSimple is a convenience method for simple single-message requests
func (c *AIClient) ChatSimple(userMessage string) (string, error) {
	messages := []AIMessage{
		{
			Role:    "user",
			Content: userMessage,
		},
	}
	return c.Chat(messages)
}
