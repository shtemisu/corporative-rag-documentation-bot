package ollama

import (
	"context"
	"corp-rag-tg-bot/internal/domain"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/client"
)

var _ domain.Embedder = (*Client)(nil)
var _ domain.Generator = (*Client)(nil)

type Client struct {
	client     *client.Client
	baseURL    string
	embedModel string
	genModel   string
}

func NewClient(baseURL string, embedModel, genModel string) *Client {
	httpClient := client.New()
	httpClient.SetTimeout(20 * time.Second)
	return &Client{
		client:     httpClient,
		baseURL:    baseURL,
		embedModel: embedModel,
		genModel:   genModel,
	}
}

type embedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type embedResponse struct {
	Model      string      `json:"model"`
	Embeddings [][]float32 `json:"embeddings"`
}

func (c *Client) Embed(ctx context.Context, text []string) ([][]float32, error) {
	req := embedRequest{Model: c.embedModel, Input: text}
	res := embedResponse{}

	url := fmt.Sprintf("%s/api/embed", c.baseURL)

	resp, err := c.client.Post(url, client.Config{
		Body: req,
		Ctx:  ctx,
	})
	if err != nil {
		return nil, fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Close()
	if resp.StatusCode() != fiber.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d: %s", resp.StatusCode(), resp.Body())
	}

	if err := resp.JSON(&res); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}
	if len(res.Embeddings) != len(text) {
		return nil, fmt.Errorf("failed to get embeddings, want: %d, but got: %d", len(res.Embeddings), len(text))
	}
	return res.Embeddings, nil
}

type generateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type generateResponse struct {
	Model    string `json:"model"`
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	req := generateRequest{
		Model:  c.genModel,
		Prompt: prompt,
		Stream: false,
	}
	res := generateResponse{}

	url := fmt.Sprintf("%s/api/generate", c.baseURL)
	resp, err := c.client.Post(url, client.Config{
		Body: req,
		Ctx:  ctx,
	})
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Close()

	if resp.StatusCode() != fiber.StatusOK {
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode(), resp.Body())
	}

	if err := resp.JSON(&res); err != nil {
		return "", fmt.Errorf("failed to parse json: %w", err)
	}
	return res.Response, nil
}
