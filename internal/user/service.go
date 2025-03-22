package user

import (
	"context"

	"archive_bot/internal/entities"
	"archive_bot/pkg/er"
	"archive_bot/pkg/logger"
)

var ErrUserNotExists error = er.New("user is not exists", "", nil)

type Repository interface {
	Save(ctx context.Context, u *User) error
	IsExists(ctx context.Context, id int64) error
	CountUsers(ctx context.Context) (int, error)
}

type service struct {
	log  *logger.Logger
	repo Repository
}

func NewService(ctx context.Context, log *logger.Logger, repo Repository) *service {
	return &service{log: log, repo: repo}
}

func (s *service) Save(ctx context.Context, event *entities.Event) error {
	u := &User{ID: event.Meta.UserID, Username: event.Meta.UserName}
	if err := s.repo.Save(ctx, u); err != nil {
		s.log.Error("failed to save user", logger.ErrAttr(err))
		return err
	}
	return nil
}

func (s *service) IsExists(ctx context.Context, event *entities.Event) error {
	if err := s.repo.IsExists(ctx, event.Meta.UserID); err != nil {
		if err == ErrUserNotFound {
			s.log.Info("new user")
			return ErrUserNotExists
		}
		s.log.Error("failed to get user", logger.ErrAttr(err))
		return err
	}

	return nil
}

func (s *service) CountUsers(ctx context.Context) (int, error) {
	return s.repo.CountUsers(ctx)
}
