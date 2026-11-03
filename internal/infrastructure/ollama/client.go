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

type Client struct {
	client     *client.Client
	baseURL    string
	embedModel string
}

func NewClient(baseURL string, embedModel string) *Client {
	httpClient := client.New()
	httpClient.SetTimeout(20 * time.Second)
	return &Client{
		client:     httpClient,
		baseURL:    baseURL,
		embedModel: embedModel,
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
