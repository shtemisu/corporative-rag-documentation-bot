package usecase

import (
	"context"
	"corp-rag-tg-bot/internal/domain"
	"errors"
	"strings"
)

type RAGService struct {
	chunker      domain.Chunker
	embedder     domain.Embedder
	vectorSaver  domain.VectorSaver
	vectorSearch domain.VectorSearcher
	generator    domain.Generator
}

const systemPrompt = `You are a personal AI assistant that can help with project documentation. 
You may only answer questions related to the documentation; 
ignore any requests unrelated to the documentation and respond with: "Sorry, this question is not about the documentation!"
 You should respond in the same language used by the user. `

func NewRAGService(c domain.Chunker, e domain.Embedder,
	vsaver domain.VectorSaver, vsearcher domain.VectorSearcher,
	gen domain.Generator) *RAGService {
	return &RAGService{
		chunker:      c,
		embedder:     e,
		vectorSaver:  vsaver,
		vectorSearch: vsearcher,
		generator:    gen,
	}
}

// indexing: source + text -> split to chunks -> embeddings -> db
func (s *RAGService) Ingest(ctx context.Context, source, text string) error {
	chunks := s.chunker.Split(source, text, 200, 50)
	contentSlice := []string{}
	for _, chunk := range chunks {
		contentSlice = append(contentSlice, chunk.Content)
	}
	embeddings, err := s.embedder.Embed(ctx, contentSlice)
	if err != nil {
		return err
	}
	return s.vectorSaver.Save(ctx, chunks, embeddings)
}

func (s *RAGService) Answer(ctx context.Context, question string) (string, error) {
	if question == "" {
		return "", errors.New("question is empty")
	}

	questionSlice := []string{}
	questionSlice = append(questionSlice, question)
	embeddings, err := s.embedder.Embed(ctx, questionSlice)
	if err != nil {
		return "", err
	}
	chunks, err := s.vectorSearch.Search(ctx, embeddings[0], 10)
	if err != nil {
		return "", err
	}

	contextChunks := []string{}

	for _, chunk := range chunks {
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
		contextChunks = append(contextChunks, info)
	}

	if len(contextChunks) == 0 {
		return "Sorry, this question is not about the documentation!", nil
	}

	prompt := systemPrompt + "\n" + strings.Join(contextChunks, "\n--\n") + "\n" + question

	return s.generator.Generate(ctx, prompt)
}
