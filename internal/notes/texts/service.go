package texts

import (
	"context"
	"archive_bot/internal/const/messages"
	"archive_bot/internal/entities"
	"archive_bot/pkg/logger"
	"time"
)

type Repository interface {
	Save(ctx context.Context, n *TextNote) (int, error)
	AllFrom(ctx context.Context, n *TextNote) ([]*TextNote, error)
	Move(ctx context.Context, n *TextNote) error
	FindLast(ctx context.Context, userID int64) (*TextNote, error)
	MoveLast(ctx context.Context, n *TextNote) error
	UpdateByID(ctx context.Context, n *TextNote) error
	RemoveByID(ctx context.Context, id int) error
}

type service struct {
	log  *logger.Logger
	repo Repository
}

func NewService(ctx context.Context, log *logger.Logger, repo Repository) *service {
	return &service{log: log, repo: repo}
}

func (s *service) Save(ctx context.Context, event *entities.Event) (int, string) {
	log := s.log.With(logger.String("operation", "texts.service.Save"))

	n := &TextNote{
		UserID:       event.Meta.UserID,
		FolderID:     event.FolderID,
		Type:         event.Type.String(),
		Description:  event.Text,
		MediaGroupID: event.MediaGroupID,
	}

	id, err := s.repo.Save(ctx, n)
	if err != nil {
		log.Error("failed to save note", logger.ErrAttr(err))
		return 0, messages.Error
	}

	if n.MediaGroupID != "" && n.Description == "" {
		return id, ""
	}

	log.Debug("result", logger.Int("id", id))
	return id, messages.NoteCreated
}

func (s *service) AllFrom(
	ctx context.Context, event *entities.Event,
) (map[int]*entities.AnswerParams, int) {
	log := s.log.With(logger.String("operation", "texts.service.AllFrom"))

	notes, err := s.repo.AllFrom(ctx, &TextNote{
		UserID: event.Meta.UserID, FolderID: event.FolderID,
	})
	if err != nil {
		if err == ErrNoTextNote {
			log.Info("notes is empty", logger.ErrAttr(err))
			return nil, 0
		}
		log.Error("failed to get all text notes", logger.ErrAttr(err))
		return nil, 0
	}

	res := make(map[int]*entities.AnswerParams, len(notes))
	for i := range notes {
		res[notes[i].ID] = &entities.AnswerParams{
			Message: notes[i].Description,
			Type:    entities.ParseType(notes[i].Type),
		}
		s.log.Debug(
			"AnswerParams",
			logger.String("Message", res[notes[i].ID].Message),
			logger.String("Type", res[notes[i].ID].Type.String()),
		)
	}

	return res, event.FolderID
}

func (s *service) Move(ctx context.Context, event *entities.Event) string {
	log := s.log.With(logger.String("operation", "texts.service.Move"))

	n := &TextNote{
		ID:       event.NoteID,
		FolderID: event.FolderID,
	}
	log.Debug("", logger.Int("note ID", n.ID), logger.Int("folder ID", n.FolderID))
	if err := s.repo.Move(ctx, n); err != nil {
		log.Error("failed to move texts note", logger.ErrAttr(err))
		return messages.Error
	}

	return messages.Moved
}

func (s *service) FindLast(ctx context.Context, event *entities.Event) (string, time.Time) {
	log := s.log.With(logger.String("operation", "texts.service.FindLast"))

	note, err := s.repo.FindLast(ctx, event.Meta.UserID)
	if err != nil {
		log.Error("failed to find texts note", logger.ErrAttr(err))
		return "", time.Time{}
	}

	return note.Description, note.CreatedAt
}

func (s *service) MoveLast(ctx context.Context, event *entities.Event) string {
	log := s.log.With(logger.String("operation", "texts.service.MoveLast"))

	n := &TextNote{
		UserID:   event.Meta.UserID,
		FolderID: event.FolderID,
	}
	if err := s.repo.MoveLast(ctx, n); err != nil {
		log.Error("failed to move text note", logger.ErrAttr(err))
		return messages.Error
	}

	return messages.Moved
}

func (s *service) UpdateByID(ctx context.Context, event *entities.Event) string {
	log := s.log.With(logger.String("operation", "texts.service.UpdateByID"))

	n := &TextNote{
		ID:          event.NoteID,
		Description: event.Text,
	}
	log.Info("", logger.Int("note ID", n.ID))
	if err := s.repo.UpdateByID(ctx, n); err != nil {
		log.Error("failed to update text note", logger.ErrAttr(err))
		return messages.Error
	}

	return n.Description
}

func (s *service) RemoveByID(ctx context.Context, id int) error {
	return s.repo.RemoveByID(ctx, id)
}
