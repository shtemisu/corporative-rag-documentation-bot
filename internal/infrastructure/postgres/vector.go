package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type VectorRepository struct {
	pool *pgxpool.Pool
}

func NewVectorRepository(p *pgxpool.Pool) *VectorRepository {
	return &VectorRepository{
		pool: p,
	}
}
