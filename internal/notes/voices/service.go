package voices

import (
	"context"

	"archive_bot/internal/entities"

	"archive_bot/pkg/logger"
)

type Repository interface {
	Save(ctx context.Context, n *Voice) (int, error)
	FindByTextsID(ctx context.Context, id int) ([]string, error)
	UpdateByTextsID(ctx context.Context, n *Voice) error
}

type service struct {
	log  *logger.Logger
	repo Repository
}

func NewService(
	ctx context.Context, log *logger.Logger, repo Repository,
) *service {
	return &service{log: log, repo: repo}
}

func (s *service) Save(ctx context.Context, event *entities.Event) int {
	log := s.log.With(logger.String("operation", "voices.service.Save"))

	id, err := s.repo.Save(ctx, &Voice{
		TextsID: event.NoteID,
		FileID:  event.FileID,
	})
	if err != nil {
		log.Error("failed to save voice note", logger.ErrAttr(err))
		return 0
	}

	return id
}

func (s *service) FindByTextsID(ctx context.Context, textsID int) []string {
	log := s.log.With(logger.String("operation", "voices.service.FindByTextsID"))

	voiceIDs, err := s.repo.FindByTextsID(ctx, textsID)
	if err != nil {
		if err == ErrNoVoice {
			log.Info("notes is empty", logger.ErrAttr(err))
			return []string{}
		}
		log.Error("failed to get all voice notes", logger.ErrAttr(err))
		return []string{}
	}

	return voiceIDs
}

func (s *service) UpdateByTextsID(ctx context.Context, event *entities.Event) error {
	log := s.log.With(logger.String("operation", "voices.service.UpdateByID"))

	n := &Voice{
		TextsID: event.NoteID,
		FileID:  event.FileID,
	}
	if err := s.repo.UpdateByTextsID(ctx, n); err != nil {
		log.Error("failed to update voice note", logger.ErrAttr(err))
		return err
	}

	return nil
}
