package usecase

import (
	"github.com/gofiber/fiber/v3/client"
)

type RAGService struct {
	LLMClient client.Client
}
