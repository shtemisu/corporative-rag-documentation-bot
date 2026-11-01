package usecase

import (
	"context"
	"corp-rag-tg-bot/internal/domain"
	"fmt"
)

type ProfileService struct {
	getter  domain.UserGetter
	updater domain.UserUpdater
	deleter domain.UserDeleter
}

func NewProfileService(g domain.UserGetter, u domain.UserUpdater, d domain.UserDeleter) *ProfileService {
	return &ProfileService{getter: g, updater: u, deleter: d}
}

func (s *ProfileService) Get(ctx context.Context, id int64) (*domain.User, error) {
	valid := validateID(id)
	if valid != nil {
		return nil, fmt.Errorf("userId is invalid: %d", id)
	}
	return s.getter.GetById(ctx, id)
}

func (s *ProfileService) Update(ctx context.Context, info *domain.User) error {
	err := validateUsername(info.Username)
	if err != nil {
		return err
	}
	return s.updater.Update(ctx, info)
}

func (s *ProfileService) Delete(ctx context.Context, id int64) error {
	return s.deleter.Delete(ctx, id)
}

func validateUsername(username string) error {
	if username == "" || len([]rune(username)) < 4 {
		return fmt.Errorf("%w: too short", domain.ErrEmptyName)
	}
	return nil
}

func validateID(id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: %d", domain.ErrInvalidData, id)
	}
	return nil
}
