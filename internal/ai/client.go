package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	APIBaseURL = "https://api.mistral.ai/v1"
)

// TODO: need to do some cleanup to remove hard-coded use of Mistral Agents

// NewClient creates a new Mistral AI client
// If agentID is provided, it will use the agent endpoint
// Otherwise, it will use the chat completions endpoint with the specified model
func NewAIClient(apiKey string, agentID string, model string) *AIClient {
	return &AIClient{
		apiKey:  apiKey,
		agentID: agentID,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Chat sends a chat completion request to the Mistral AI API
// If agentID is set, it uses the agent endpoint, otherwise uses chat completions
func (c *AIClient) Chat(messages []AIMessage) (string, error) {
	return c.chatWithAgent(messages)
}

// chatWithAgent sends a request to the agent completion endpoint
func (c *AIClient) chatWithAgent(messages []AIMessage) (string, error) {
	reqBody := AgentRequest{
		AgentID:  c.agentID,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", APIBaseURL+"/agents/completions", bytes.NewBuffer(jsonData))
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
		return "", fmt.Errorf("API error: %s", errResp.Error.Message)
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
