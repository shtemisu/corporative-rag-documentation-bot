package usecase

import "corp-rag-tg-bot/internal/domain"

type RAGService struct {
	embedder domain.Embedder
}
