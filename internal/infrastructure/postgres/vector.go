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

	source := chunks[0].Source
	for _, c := range chunks {
		if c.Source != source {
			return fmt.Errorf("mixed sources in one Save: %q and %q", source, c.Source)
		}
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return mapPgErr(err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM chunks WHERE source = $1`, source); err != nil {
		return mapPgErr(err)
	}

	batch := &pgx.Batch{}
	query := `
		INSERT INTO chunks (content, source, embedding, chunks_index, heading) 
		VALUES ($1, $2, $3, $4, $5)
	`
	for i, chunk := range chunks {
		batch.Queue(query, chunk.Content, chunk.Source, pgvector.NewVector(embeddings[i]), chunk.Index, chunk.Heading)
	}

	br := tx.SendBatch(ctx, batch)
	for i := range chunks {
		if _, err := br.Exec(); err != nil {
			_ = br.Close()
			return mapPgErr(fmt.Errorf("batch item %d: %w", i, err))
		}
	}
	if err := br.Close(); err != nil {
		return mapPgErr(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return mapPgErr(err)
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
		SELECT content, source, chunks_index, heading, (embedding <=> $1)::float8 AS distance
		FROM chunks
		ORDER BY embedding <=> $1 ASC
		LIMIT $2
	`

	vec := pgvector.NewVector(embeddings)

	rows, err := r.pool.Query(ctx, query, vec, topK)
	if err != nil {
		return nil, mapPgErr(err)
	}
	defer rows.Close()
	var chunks []domain.Chunk

	for rows.Next() {
		var chunk domain.Chunk
		err := rows.Scan(
			&chunk.Content,
			&chunk.Source,
			&chunk.Index,
			&chunk.Heading,
			&chunk.Distance,
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
