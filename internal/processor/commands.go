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
	return messages.StartCommand, buttons.Folders
}

func (p *processor) Folders(ctx context.Context, event *entities.Event) map[string]string {
	// log := p.log.With(logger.String("operation", "processor.Folders"))
	p.fm.SetCurrentFolderID(
		event.Meta.UserID,
		p.fm.service.DefaultFolderID(ctx, event.Meta.UserID),
	)

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
	if createdAt.IsZero() {
		log.Error(
			"failed to find last note",
			logger.String("event", event.String()),
		)
		return messages.Error
	}
	if event.Meta.Date.Sub(createdAt) < 60*time.Second {
		event.FolderID = folderID
		return p.nm.texts.MoveLast(ctx, event)
	}

	p.fm.SetCurrentFolderID(event.Meta.UserID, folderID)
	return ""
}
