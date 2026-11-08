package postgres

import (
	"context"
	"corp-rag-tg-bot/internal/domain"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

var _ domain.VectorSaver = (*VectorRepository)(nil)

type VectorRepository struct {
	pool *pgxpool.Pool
}

func NewVectorRepository(p *pgxpool.Pool) *VectorRepository {
	return &VectorRepository{
		pool: p,
	}
}

func (r *VectorRepository) Save(ctx context.Context, chunks []domain.Chunk, embeddings [][]float32) error {
	if len(chunks) != len(embeddings) {
		return fmt.Errorf("chunks count (%d) does not match embeddings count (%d)", len(chunks), len(embeddings))
	}
	if len(chunks) == 0 {
		return nil
	}

	batch := &pgx.Batch{}

	query := `
		INSERT INTO chunks (content, source, embedding, chunks_index, heading) 
		VALUES ($1, $2, $3, $4, $5)
	`

	for idx, chunk := range chunks {
		vec := pgvector.NewVector(embeddings[idx])
		batch.Queue(query, chunk.Content, chunk.Source, vec, chunk.Index, chunk.Heading)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()

	for i := range chunks {
		_, err := br.Exec()
		if err != nil {
			return mapPgErr(fmt.Errorf("failed to execute batch item %d: %w", i, err))
		}
	}
	return nil
}

func (r *VectorRepository) Search(ctx context.Context, embeddings []float32, topK int) ([]domain.Chunk, error) {
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("embeddings count is equal zero: %d", len(embeddings))
	}

	if topK <= 0 {
		topK = 5
	}

	query := `
		SELECT content, source, chunks_index, heading FROM chunks
		ORDER BY embedding <=> $1 ASC
		LIMIT $2
	`

	vec := pgvector.NewVector(embeddings)

	rows, err := r.pool.Query(ctx, query, vec, topK)
	if err != nil {
		return nil, mapPgErr(err)
	}

	var chunks []domain.Chunk

	for rows.Next() {
		var chunk domain.Chunk
		err := rows.Scan(
			&chunk.Content,
			&chunk.Source,
			&chunk.Index,
			&chunk.Heading,
		)

		if err != nil {
			return nil, mapPgErr(err)
		}

		chunks = append(chunks, chunk)
	}

	if err := rows.Err(); err != nil {
		return nil, mapPgErr(err)
	}

	return chunks, nil
}
