package processor

import (
	"cmp"
	"context"
	"time"

	"archive_bot/internal/const/buttons"
	"archive_bot/internal/const/messages"
	"archive_bot/internal/entities"
	"archive_bot/pkg/logger"
)

func (p *processor) Start(ctx context.Context, event *entities.Event) (string, string) {
	log := p.log.With(logger.String("operation", "processor.Start"))

	if p.fm.service.DefaultFolderID(ctx, event.Meta.UserID) == 0 {
		if err := p.fm.service.SaveDefault(ctx, event); err != nil {
			log.Error("save default folder error", logger.ErrAttr(err))
		}
	}

	return messages.StartCommand, buttons.Folders
}

func (p *processor) Folders(ctx context.Context, event *entities.Event) map[string]string {
	// log := p.log.With(logger.String("operation", "processor.Folders"))
	p.fm.SetCurrentFolderID(event.Meta.UserID, event.FolderID)

	return p.fm.service.All(ctx, event)
}

func (p *processor) Save(ctx context.Context, event *entities.Event) *entities.AnswerParams {
	// log := p.log.With(logger.String("operation", "processor.Save"))

	event.FolderID = cmp.Or(
		p.fm.CurrentFolderID(event.Meta.UserID),
		p.fm.service.DefaultFolderID(ctx, event.Meta.UserID),
	)

	noteID, message := p.nm.texts.Save(ctx, event)
	event.NoteID = noteID
	ap := entities.AnswerParams{Message: message}
	switch event.Type {
	case entities.Photo:
		p.nm.photos.Save(ctx, event)
	case entities.Document:
		p.nm.documents.Save(ctx, event)
	case entities.Video:
		p.nm.videos.Save(ctx, event)
	case entities.Audio:
		p.nm.audios.Save(ctx, event)
	case entities.Animation:
		p.nm.ani.Save(ctx, event)
	case entities.Voice:
		p.nm.voices.Save(ctx, event)
	}

	return &ap
}

func (p *processor) SaveTo(ctx context.Context, event *entities.Event) string {
	log := p.log.With(logger.String("operation", "processor.SaveTo"))
	folderID, err := p.fm.service.FindOrCreate(ctx, event)
	if err != nil {
		log.Error(
			"failed to find or create folder",
			logger.String("event", event.String()),
			logger.ErrAttr(err),
		)
		return messages.Error
	}

	_, createdAt := p.nm.texts.FindLast(ctx, event)
	log.Debug(
		createdAt.String(),
		logger.Bool("createdAt", createdAt.IsZero()),
	)
	if createdAt.IsZero() {
		log.Debug(
			"last note is not exists",
			logger.String("event", event.String()),
		)
	} else if event.Meta.Date.Sub(createdAt) < 1*time.Second {
		log.Debug(
			"last note is exists",
			logger.String("event", event.String()),
		)
		event.FolderID = folderID
		return p.nm.texts.MoveLast(ctx, event)
	}

	p.fm.SetCurrentFolderID(event.Meta.UserID, folderID)
	event.FolderID = folderID
	log.Debug(
		"SaveTo",
		logger.Int("event.FolderID", event.FolderID),
		logger.Int("CurrentFolderID", p.fm.CurrentFolderID(event.Meta.UserID)),
		logger.Bool("Meta.Date - createdAt", event.Meta.Date.Sub(createdAt) < 1*time.Second),
	)

	return messages.Moved
}
