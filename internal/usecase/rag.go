package usecase

import (
	"context"
	"corp-rag-tg-bot/internal/domain"
	"errors"
	"log"
	"strings"
)

type RAGService struct {
	chunker      domain.Chunker
	embedder     domain.Embedder
	vectorSaver  domain.VectorSaver
	vectorSearch domain.VectorSearcher
	generator    domain.Generator
	maxDistance  float64
}

const (
	DefaultMaxDistance = 0.5

	searchTopK = 10
)

const systemPrompt = `You are a personal AI assistant that can help with project documentation.
You may only answer questions related to the documentation;
ignore any requests unrelated to the documentation and respond with: "Sorry, this question is not about the documentation!"
You should respond in the same language used by the user. `

const notAboutDocumentation = "Sorry, this question is not about the documentation!"

func NewRAGService(c domain.Chunker, e domain.Embedder,
	vsaver domain.VectorSaver, vsearcher domain.VectorSearcher,
	gen domain.Generator, maxDistance float64) *RAGService {
	if maxDistance <= 0 {
		maxDistance = DefaultMaxDistance
	}
	return &RAGService{
		chunker:      c,
		embedder:     e,
		vectorSaver:  vsaver,
		vectorSearch: vsearcher,
		generator:    gen,
		maxDistance:  maxDistance,
	}
}

// indexing: source + text -> split to chunks -> embeddings -> db
func (s *RAGService) Ingest(ctx context.Context, source, text string) error {
	chunks := s.chunker.Split(source, text, 200, 50)
	if len(chunks) == 0 {
		return nil
	}
	contentSlice := []string{}
	for _, chunk := range chunks {
		contentSlice = append(contentSlice, chunk.Content)
	}
	embeddings, err := s.embedder.Embed(ctx, contentSlice)
	if len(embeddings) == 0 {
		return nil
	}
	if err != nil {
		return err
	}
	return s.vectorSaver.Save(ctx, chunks, embeddings)
}

type chunkKey struct {
	source string
	index  int
}

func (s *RAGService) Answer(ctx context.Context, question string) (string, error) {
	if question == "" {
		return "", errors.New("question is empty")
	}

	embeddings, err := s.embedder.Embed(ctx, []string{question})
	if err != nil {
		return "", err
	}
	if len(embeddings) == 0 {
		return "", errors.New("embedder returned no embeddings")
	}

	chunks, err := s.vectorSearch.Search(ctx, embeddings[0], searchTopK)
	if err != nil {
		return "", err
	}

	contextChunks := make([]string, 0, len(chunks))
	seen := make(map[chunkKey]struct{}, len(chunks))
	bestDistance := -1.0
	if len(chunks) > 0 {
		bestDistance = chunks[0].Distance
	}

	for _, chunk := range chunks {
		if chunk.Distance > s.maxDistance {
			continue
		}
		key := chunkKey{source: chunk.Source, index: chunk.Index}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		info := ""
		if len(chunk.Source) != 0 {
			info += "Source: " + chunk.Source + "\n"
		}
		if len(chunk.Heading) != 0 {
			info += "Heading: " + chunk.Heading + "\n"
		}
		if len(chunk.Content) != 0 {
			info += "Content: " + chunk.Content + "\n"
		}
		if info != "" {
			contextChunks = append(contextChunks, info)
		}
	}

	if len(contextChunks) == 0 {
		log.Printf("rag: question rejected, best_distance=%.3f max_distance=%.3f",
			bestDistance, s.maxDistance)
		return notAboutDocumentation, nil
	}

	log.Printf("rag: context built, best_distance=%.3f max_distance=%.3f kept=%d retrieved=%d",
		bestDistance, s.maxDistance, len(contextChunks), len(chunks))

	prompt := systemPrompt + "\n\nDocumentation:\n" + strings.Join(contextChunks, "\n---\n") +
		"\n\nQuestion:\n" + question + "\n\nAnswer:"

	return s.generator.Generate(ctx, prompt)
}
