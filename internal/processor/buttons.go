package processor

import (
	"context"
	"strconv"
	"strings"

	"archive_bot/internal/const/messages"
	"archive_bot/internal/entities"

	"archive_bot/pkg/logger"
)

func (p *processor) SelectFolder(
	ctx context.Context,
	event *entities.Event,
) (map[int]*entities.AnswerParams, string) {
	folderID, err := strconv.Atoi(strings.Split(event.Text, "_")[1])
	if err != nil {
		return nil, ""
	}
	event.FolderID = folderID
	p.fm.SetCurrentFolderID(event.Meta.UserID, folderID)

	folderName, err := p.fm.service.Find(ctx, event)
	if err != nil {
		return nil, ""
	}
	answerParamsMap, _ := p.nm.texts.AllFrom(ctx, event)

	for textsID, ap := range answerParamsMap {
		if ap.Message == "" {
			ap.Message = messages.EmptyMessage
		}
		switch ap.Type {
		case entities.Photo:
			ap.FileIDs = p.nm.photos.FindByTextsID(ctx, textsID)
		case entities.Document:
			ap.FileIDs = p.nm.documents.FindByTextsID(ctx, textsID)
		case entities.Video:
			ap.FileIDs = p.nm.videos.FindByTextsID(ctx, textsID)
		case entities.Audio:
			ap.FileIDs = p.nm.audios.FindByTextsID(ctx, textsID)
		case entities.Animation:
			ap.FileIDs = p.nm.ani.FindByTextsID(ctx, textsID)
		case entities.Voice:
			ap.FileIDs = p.nm.voices.FindByTextsID(ctx, textsID)
		}
	}

	return answerParamsMap, folderName
}

func (p *processor) RemoveNote(ctx context.Context, event *entities.Event) string {
	log := logger.L(ctx).With(logger.String("operation", "processor.RemoveNote"))

	if err := p.nm.texts.RemoveByID(ctx, event.NoteID); err != nil {
		log.Error("", logger.ErrAttr(err))
		return messages.Error
	}
	return messages.NoteRemoved
}

func (p *processor) AddFolderStart(ctx context.Context, event *entities.Event) string {
	log := logger.L(ctx).With(logger.String("operation", "processor.AddFolderStart"))

	state := p.fm.setStateCreate(event.Meta.UserID)
	if state.MessageID == 0 {
		state.MessageID = event.Meta.MessageID
	}
	switch state.FSM.Current() {
	case StartCreate:
		if err := state.FSM.Event(ctx, "begin"); err != nil {
			log.Error("failed to transit state", logger.ErrAttr(err))
			return messages.Error
		}

		return messages.AskFolderName
	default:
		return messages.Error
	}

}

func (p *processor) AddFolderEnd(ctx context.Context, event *entities.Event) string {
	log := logger.L(ctx).With(logger.String("operation", "processor.AddFolderEnd"))
	state := p.fm.stateCreate(event.Meta.UserID)
	if state == nil {
		return ""
	}
	switch state.FSM.Current() {
	case SelectCreate:
		if err := state.FSM.Event(ctx, "provide_name"); err != nil {
			log.Error("failed to transit state", logger.ErrAttr(err))
			return messages.Error
		}

		event.Meta.MessageID = state.MessageID
		state.MessageID = 0

		return p.fm.service.Save(ctx, event)
	default:
		return ""
	}
}

func (p *processor) DeleteFolderStart(ctx context.Context, event *entities.Event) string {
	log := logger.L(ctx).With(logger.String("operation", "processor.DeleteFolderStart"))

	state := p.fm.setStateDelete(event.Meta.UserID)

	switch state.FSM.Current() {
	case StartDelete:
		if err := state.FSM.Event(ctx, "begin"); err != nil {
			log.Error("failed to transit state", logger.ErrAttr(err))
			return messages.Error
		}

		return messages.ChooseFolderToDelete
	default:
		return ""
	}

}

func (p *processor) DeleteFolderEnd(ctx context.Context, event *entities.Event) string {
	log := logger.L(ctx).With(logger.String("operation", "processor.DeleteFolderEnd"))

	state := p.fm.stateDelete(event.Meta.UserID)
	if state == nil {
		return ""
	}

	switch state.FSM.Current() {
	case SelectDelete:
		if err := state.FSM.Event(ctx, "provide_name"); err != nil {
			log.Error("failed to transit state", logger.ErrAttr(err))
			return messages.Error
		}

		id, err := strconv.Atoi(strings.Split(event.Text, "_")[1])
		if err != nil {
			log.Error("failed to transit state", logger.ErrAttr(err))
			return messages.FolderNotExists
		}

		if id == p.fm.service.DefaultFolderID(ctx, event.Meta.UserID) {
			return messages.WrongFolder
		}

		if err := p.fm.service.RemoveByID(ctx, id); err != nil {
			log.Error("failed to transit state", logger.ErrAttr(err))
			return messages.FolderNotExists
		}

		return messages.FolderDeleted
	default:
		return ""
	}
}

func (p *processor) MoveNoteStart(ctx context.Context, event *entities.Event) string {
	log := logger.L(ctx).With(logger.String("operation", "processor.MoveNoteStart"))

	log.Debug("",
		logger.Int64("UserID", event.Meta.UserID),
		logger.Int("FolderID", event.FolderID),
		logger.Int("NoteID", event.NoteID),
	)
	state := p.nm.MoveState(event.Meta.UserID, event)
	if state.ParentFolderID == 0 && state.NoteID == 0 {
		state.ParentFolderID = event.FolderID
		state.NoteID = event.NoteID
	}

	switch state.FSM.Current() {
	case StartMove:
		if err := state.FSM.Event(ctx, "begin"); err != nil {
			log.Error("failed to transit state", logger.ErrAttr(err))
			return messages.Error
		}

		return messages.ChooseFolderToMove
	default:
		return ""
	}
}

func (p *processor) MoveNoteEnd(ctx context.Context, event *entities.Event) string {
	log := logger.L(ctx).With(logger.String("operation", "processor.MoveNoteEnd"))

	state := p.nm.MoveState(event.Meta.UserID, event)

	switch state.FSM.Current() {
	case SelectMove:
		if err := state.FSM.Event(ctx, "provide_ID"); err != nil {
			log.Error("failed to transit state", logger.ErrAttr(err))
			return messages.Error
		}

		log.Debug("",
			logger.Int64("UserID", event.Meta.UserID),
			logger.Int("FolderID", event.FolderID),
			logger.Int("ParentFolderID", state.ParentFolderID),
			logger.Int("NoteID", state.NoteID),
		)

		event.NoteID = state.NoteID
		event.FolderID, _ = strconv.Atoi(strings.Split(event.Text, "_")[1])
		message := p.nm.texts.Move(ctx, event)
		event.FolderID = state.ParentFolderID

		state.ParentFolderID = 0
		state.NoteID = 0

		return message
	default:
		return ""
	}
}
