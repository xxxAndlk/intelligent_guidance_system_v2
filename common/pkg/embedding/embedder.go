package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrEmptyText      = errors.New("text cannot be empty")
	ErrEmptyBatch     = errors.New("batch cannot be empty")
	ErrAPITimeout     = errors.New("embedding API timeout")
	ErrAPIError       = errors.New("embedding API error")
	ErrInvalidResponse = errors.New("invalid API response")
)

type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
	Dimension() int
	ModelName() string
}

type SiliconFlowEmbedder struct {
	apiKey  string
	apiURL  string
	model   string
	timeout time.Duration
	client  *http.Client
}

type SiliconFlowConfig struct {
	APIKey  string        `json:"api_key" yaml:"api_key"`
	APIURL  string        `json:"api_url" yaml:"api_url"`
	Model   string        `json:"model" yaml:"model"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
}

func NewSiliconFlowEmbedder(cfg SiliconFlowConfig) (*SiliconFlowEmbedder, error) {
	if cfg.APIKey == "" {
		return nil, errors.New("API key is required")
	}
	if cfg.APIURL == "" {
		cfg.APIURL = "https://api.siliconflow.cn/v1/embeddings"
	}
	if cfg.Model == "" {
		cfg.Model = "BAAI/bge-large-zh-v1.5"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	return &SiliconFlowEmbedder{
		apiKey:  cfg.APIKey,
		apiURL:  cfg.APIURL,
		model:   cfg.Model,
		timeout: cfg.Timeout,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}, nil
}

type embeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type embeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

func (e *SiliconFlowEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, ErrEmptyText
	}

	embeddings, err := e.EmbedBatch(ctx, []string{text})
	if err != nil {
		return nil, err
	}

	if len(embeddings) == 0 {
		return nil, ErrInvalidResponse
	}

	return embeddings[0], nil
}

func (e *SiliconFlowEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, ErrEmptyBatch
	}

	reqBody := embeddingRequest{
		Model: e.model,
		Input: texts,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", e.apiURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.apiKey)

	resp, err := e.client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrAPITimeout
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var embResp embeddingResponse
	if err := json.Unmarshal(body, &embResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if embResp.Error != nil {
		return nil, fmt.Errorf("%w: %s", ErrAPIError, embResp.Error.Message)
	}

	if len(embResp.Data) == 0 {
		return nil, ErrInvalidResponse
	}

	embeddings := make([][]float32, len(embResp.Data))
	for i, d := range embResp.Data {
		embeddings[i] = d.Embedding
	}

	return embeddings, nil
}

func (e *SiliconFlowEmbedder) Dimension() int {
	switch e.model {
	case "BAAI/bge-large-zh-v1.5", "BAAI/bge-large-en-v1.5":
		return 1024
	case "BAAI/bge-base-en-v1.5", "BAAI/bge-base-zh-v1.5":
		return 768
	case "BAAI/bge-small-en-v1.5", "BAAI/bge-small-zh-v1.5":
		return 512
	default:
		return 1024
	}
}

func (e *SiliconFlowEmbedder) ModelName() string {
	return e.model
}

type MockEmbedder struct {
	dimension int
	modelName string
}

func NewMockEmbedder(dimension int, modelName string) *MockEmbedder {
	if dimension <= 0 {
		dimension = 768
	}
	if modelName == "" {
		modelName = "mock-embedder"
	}
	return &MockEmbedder{
		dimension: dimension,
		modelName: modelName,
	}
}

func (e *MockEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, ErrEmptyText
	}
	embedding := make([]float32, e.dimension)
	for i := range embedding {
		embedding[i] = 0.1
	}
	return embedding, nil
}

func (e *MockEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, ErrEmptyBatch
	}
	embeddings := make([][]float32, len(texts))
	for i := range texts {
		embeddings[i] = make([]float32, e.dimension)
		for j := range embeddings[i] {
			embeddings[i][j] = 0.1
		}
	}
	return embeddings, nil
}

func (e *MockEmbedder) Dimension() int {
	return e.dimension
}

func (e *MockEmbedder) ModelName() string {
	return e.modelName
}