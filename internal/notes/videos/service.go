package videos

import (
	"context"

	"archive_bot/internal/entities"

	"archive_bot/pkg/logger"
)

type Repository interface {
	Save(ctx context.Context, n *Video) (int, error)
	FindByTextsID(ctx context.Context, id int) ([]string, error)
	UpdateByTextsID(ctx context.Context, n *Video) error
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
	log := s.log.With(logger.String("operation", "videos.service.Save"))

	id, err := s.repo.Save(ctx, &Video{
		TextsID:      event.NoteID,
		FileID:       event.FileID,
		MediaGroupID: event.MediaGroupID,
	})
	if err != nil {
		log.Error("failed to save video note", logger.ErrAttr(err))
		return 0
	}

	return id
}

func (s *service) FindByTextsID(ctx context.Context, textsID int) []string {
	log := s.log.With(logger.String("operation", "videos.service.FindByTextsID"))

	videoIDs, err := s.repo.FindByTextsID(ctx, textsID)
	if err != nil {
		if err == ErrNoVideo {
			log.Info("notes is empty", logger.ErrAttr(err))
			return []string{}
		}
		log.Error("failed to get all video notes", logger.ErrAttr(err))
		return []string{}
	}

	return videoIDs
}

func (s *service) UpdateByTextsID(ctx context.Context, event *entities.Event) error {
	log := s.log.With(logger.String("operation", "videos.service.UpdateByID"))

	n := &Video{
		TextsID:      event.NoteID,
		FileID:       event.FileID,
		MediaGroupID: event.MediaGroupID,
	}
	if err := s.repo.UpdateByTextsID(ctx, n); err != nil {
		log.Error("failed to update video note", logger.ErrAttr(err))
		return err
	}

	return nil
}
