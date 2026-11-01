package usecase

import (
	"context"
	"corp-rag-tg-bot/internal/domain"
	"time"
)

type RegistrationService struct {
	exists domain.UserChecker
	saver  domain.UserSaver
	clock  func() time.Time
}

func NewRegistrationService(e domain.UserChecker, s domain.UserSaver, c func() time.Time) *RegistrationService {
	return &RegistrationService{
		exists: e,
		saver:  s,
		clock:  c,
	}
}

func (s *RegistrationService) Register(ctx context.Context, id int64, username string) error {
	exists, err := s.exists.Exists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrUserAlreadyExists
	}

	u, err := domain.NewUser(id, username, s.clock())
	if err != nil {
		return err
	}
	return s.saver.Create(ctx, u)
}
