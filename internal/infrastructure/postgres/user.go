package postgres

import (
	"context"
	"corp-rag-tg-bot/internal/domain"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	_ domain.UserSaver   = (*UserRepository)(nil)
	_ domain.UserChecker = (*UserRepository)(nil)
	_ domain.UserGetter  = (*UserRepository)(nil)
	_ domain.UserUpdater = (*UserRepository)(nil)
	_ domain.UserDeleter = (*UserRepository)(nil)
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(p *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: p}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	q := `INSERT INTO users (id, username, created_at) VALUES ($1, $2, $3)`

	if _, err := r.pool.Exec(ctx, q, user.ID, user.Username, user.CreatedAt); err != nil {
		return mapPgErr(err)
	}
	return nil
}

func (r *UserRepository) GetById(ctx context.Context, id int64) (*domain.User, error) {
	q := `SELECT id, username, created_at FROM users WHERE id = $1`
	user := domain.User{}
	err := r.pool.QueryRow(ctx, q, id).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, mapPgErr(err)
	}

	return &user, nil
}

func (r *UserRepository) Exists(ctx context.Context, id int64) (bool, error) {
	q := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	var ok bool
	if err := r.pool.QueryRow(ctx, q, id).Scan(&ok); err != nil {
		return false, mapPgErr(err)
	}
	return ok, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	q := `UPDATE users SET username = $1 where id = $2`
	tag, err := r.pool.Exec(ctx, q, user.Username, user.ID)
	if err != nil {
		return mapPgErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	q := `DELETE FROM users WHERE id = $1`
	tag, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return mapPgErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
