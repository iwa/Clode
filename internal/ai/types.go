package ai

import "net/http"

// AIMessage represents a chat message in AI thread
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIClient represents a Mistral AI API client
type AIClient struct {
	apiKey     string
	httpClient *http.Client
	agentID    string
}

// ChatRequest represents a request to the chat completion endpoint
type ChatRequest struct {
	Model       string      `json:"model,omitempty"`
	Messages    []AIMessage `json:"messages"`
	Temperature float64     `json:"temperature,omitempty"`
	MaxTokens   int         `json:"max_tokens,omitempty"`
}

// AgentRequest represents a request to the agent completion endpoint
type AgentRequest struct {
	AgentID  string      `json:"agent_id"`
	Messages []AIMessage `json:"messages"`
}

// ChatResponse represents a response from the chat completion endpoint
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}
