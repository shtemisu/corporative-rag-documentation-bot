package usecase

import (
	"context"
	"corp-rag-tg-bot/internal/domain"
)

type RAGService struct {
	chunker      domain.Chunker
	embedder     domain.Embedder
	vectorSaver  domain.VectorSaver
	vectorSearch domain.VectorSearcher
	generator    domain.Generator
}

const prompt = `You are a personal AI assistant that can help with project documentation. 
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
func (s *RAGService) Injest(ctx context.Context, source, text string) error {
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
	// TODO:
	// 1. embeddings, err := s.embedder.Embed(ctx, []string{question})
	// 2. chunks, err := s.vectorSearch.Search(ctx, embeddings, topK)
	// 3. contextStr := склеить Content из chunks через "\n---\n"
	// 4. prompt := сформурировать системное сообщение + вопрос + контекст
	// 5. return s.generator.Generate(ctx, prompt)
	questionSlice := []string{}
	questionSlice = append(questionSlice, question)
	embeddings, err := s.embedder.Embed(ctx, questionSlice)
	if err != nil {
		return "", err
	}
	chunks, err := s.vectorSearch.Search(ctx, embeddings, 10)
	if err != nil {
		return "", err
	}

	contextStr := []string{}
	for _, chunk := range chunks {
		contextStr = append(contextStr, chunk.Content)
	}

	userPromptJoined := question + "\n" + prompt

	return s.generator.Generate(ctx, userPromptJoined)
}
