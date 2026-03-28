// Package llm provides LLM client implementations for AI inference.
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"intelligent_guidance_system_v2/service/ai/internal/domain/repository"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrLLMConnection     = errors.New("llm connection error")
	ErrLLMRequestFailed  = errors.New("llm request failed")
	ErrLLMResponseInvalid = errors.New("llm response invalid")
	ErrEmbeddingFailed   = errors.New("embedding generation failed")
)

type DeepSeekConfig struct {
	BaseURL      string
	APIKey       string
	Model        string
	EmbedModel   string
	MaxTokens    int
	Temperature  float32
	Timeout      time.Duration
}

type DeepSeekClient struct {
	config     DeepSeekConfig
	httpClient *http.Client
	logger     *log.Helper
}

func NewDeepSeekClient(cfg DeepSeekConfig, logger log.Logger) (*DeepSeekClient, error) {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.deepseek.com"
	}
	if cfg.Model == "" {
		cfg.Model = "deepseek-chat"
	}
	if cfg.EmbedModel == "" {
		cfg.EmbedModel = "deepseek-embedding"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	
	return &DeepSeekClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		logger: log.NewHelper(logger),
	}, nil
}

type ChatRequest struct {
	Model       string          `json:"model"`
	Messages    []ChatMessage   `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float32         `json:"temperature,omitempty"`
	Stream      bool            `json:"stream"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []ChatChoice   `json:"choices"`
	Usage   ChatUsage      `json:"usage"`
}

type ChatChoice struct {
	Index        int          `json:"index"`
	Message      ChatMessage  `json:"message"`
	FinishReason string       `json:"finish_reason"`
}

type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type EmbedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type EmbedResponse struct {
	Object  string        `json:"object"`
	Data    []EmbedData   `json:"data"`
	Model   string        `json:"model"`
	Usage   EmbedUsage    `json:"usage"`
}

type EmbedData struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float32 `json:"embedding"`
}

type EmbedUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

func (c *DeepSeekClient) Chat(ctx context.Context, messages []ChatMessage) (string, error) {
	req := ChatRequest{
		Model:       c.config.Model,
		Messages:    messages,
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
		Stream:      false,
	}
	
	resp, err := c.doRequest(ctx, "/v1/chat/completions", req)
	if err != nil {
		return "", err
	}
	
	var chatResp ChatResponse
	if err := json.Unmarshal(resp, &chatResp); err != nil {
		return "", fmt.Errorf("%w: %v", ErrLLMResponseInvalid, err)
	}
	
	if len(chatResp.Choices) == 0 {
		return "", ErrLLMResponseInvalid
	}
	
	return chatResp.Choices[0].Message.Content, nil
}

func (c *DeepSeekClient) ChatWithSystem(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userMessage},
	}
	return c.Chat(ctx, messages)
}

func (c *DeepSeekClient) doRequest(ctx context.Context, path string, body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to marshal request", ErrLLMRequestFailed)
	}
	
	url := c.config.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLLMConnection, err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLLMRequestFailed, err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.logger.Errorf("LLM API error: status=%d, body=%s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("%w: status code %d", ErrLLMRequestFailed, resp.StatusCode)
	}
	
	return io.ReadAll(resp.Body)
}

func (c *DeepSeekClient) Embed(ctx context.Context, text string) ([]float32, error) {
	embeddings, err := c.EmbedBatch(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	
	if len(embeddings) == 0 {
		return nil, ErrEmbeddingFailed
	}
	
	return embeddings[0], nil
}

func (c *DeepSeekClient) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	req := EmbedRequest{
		Model: c.config.EmbedModel,
		Input: texts,
	}
	
	resp, err := c.doRequest(ctx, "/v1/embeddings", req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrEmbeddingFailed, err)
	}
	
	var embedResp EmbedResponse
	if err := json.Unmarshal(resp, &embedResp); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLLMResponseInvalid, err)
	}
	
	embeddings := make([][]float32, len(embedResp.Data))
	for i, d := range embedResp.Data {
		embeddings[i] = d.Embedding
	}
	
	return embeddings, nil
}

type DiagnosisPromptBuilder struct{}

func NewDiagnosisPromptBuilder() *DiagnosisPromptBuilder {
	return &DiagnosisPromptBuilder{}
}

func (b *DiagnosisPromptBuilder) BuildSystemPrompt() string {
	return `You are an AI medical diagnosis assistant. Your role is to:
1. Analyze patient symptoms and medical history
2. Suggest potential diseases based on symptom patterns
3. Recommend appropriate medical specialties and doctors
4. Suggest relevant medications (with proper disclaimers)

Important guidelines:
- Always recommend consulting a real doctor for definitive diagnosis
- Provide confidence levels for your recommendations
- Consider symptom combinations and severity
- Ask clarifying questions when symptoms are ambiguous`
}

func (b *DiagnosisPromptBuilder) BuildAnalysisPrompt(symptoms, conversationHistory string) string {
	return fmt.Sprintf(`Please analyze the following patient symptoms and provide recommendations.

Symptoms:
%s

Conversation History:
%s

Provide your analysis in JSON format with:
- potential_diseases: list of possible conditions with confidence scores
- recommended_specialties: medical specialties to consult
- recommended_tests: diagnostic tests that may be helpful
- follow_up_questions: any clarifying questions needed
- urgency_level: low/medium/high based on symptom severity`, symptoms, conversationHistory)
}

func (b *DiagnosisPromptBuilder) BuildSymptomCollectionPrompt(existingSymptoms string) string {
	return fmt.Sprintf(`You are collecting symptoms from a patient.

Current symptoms recorded:
%s

Ask follow-up questions to gather more details about:
- Duration and onset of symptoms
- Severity changes
- Related symptoms
- Medical history
- Current medications

Be empathetic and clear. Ask one question at a time.`, existingSymptoms)
}

type EmbeddingRepoAdapter struct {
	client *DeepSeekClient
}

func NewEmbeddingRepoAdapter(client *DeepSeekClient) repository.EmbeddingRepository {
	return &EmbeddingRepoAdapter{client: client}
}

func (a *EmbeddingRepoAdapter) Embed(ctx context.Context, text string) ([]float32, error) {
	return a.client.Embed(ctx, text)
}

func (a *EmbeddingRepoAdapter) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	return a.client.EmbedBatch(ctx, texts)
}