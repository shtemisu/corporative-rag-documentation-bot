package domain

import (
	"context"
	"time"
)

type User struct {
	ID        int64     // chatID in Telegram
	Username  string    // custom username in telegram
	CreatedAt time.Time // log-time in db
}

func NewUser(id int64, username string, clock time.Time) (*User, error) {
	if username == "" || id <= 0 {
		return nil, ErrInvalidData
	}

	return &User{ID: id, Username: username, CreatedAt: clock.UTC()}, nil
}

type UserSaver interface {
	Create(ctx context.Context, info *User) error
}

type UserGetter interface {
	GetById(ctx context.Context, id int64) (*User, error)
}

type UserChecker interface {
	Exists(ctx context.Context, id int64) (bool, error)
}

type UserUpdater interface {
	Update(ctx context.Context, info *User) error
}

type UserDeleter interface {
	Delete(ctx context.Context, ID int64) error
}
