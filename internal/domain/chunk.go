package domain

import "context"

type Chunk struct {
	Source  string
	Heading string
	Index   int
	Content string
}

type VectorSaver interface {
	Save(ctx context.Context, chunks []Chunk, embeddings [][]float32) error
}

type VectorSearcher interface {
	Search(ctx context.Context, embedding [][]float32, topK int) ([]Chunk, error)
}

type Generator interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type Chunker interface {
	Split(source, text string, maxLen, overLap int) []Chunk
}
