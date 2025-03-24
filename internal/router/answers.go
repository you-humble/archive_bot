package router

import (
	"context"
	"sort"
	"strconv"

	"archive_bot/internal/const/buttons"
	"archive_bot/internal/const/messages"
	"archive_bot/internal/entities"
	"archive_bot/pkg/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (r *router) doDefaultCallback(ctx context.Context, b *bot.Bot, event *entities.Event) {
	log := logger.L(ctx).With(logger.String("operation", "router.doDefaultCallback"))
	if message := r.process.DeleteFolderEnd(ctx, event); message != "" {
		btns := r.process.Folders(ctx, event)
		event.IsEdited = true
		go func() {
			r.deleteMessages(ctx, b, event)
			if message != messages.WrongFolder {
				r.sendAnswers(ctx, b, []*entities.Answer{sendFoldersList(event, btns, false)})
				r.process.SetInt(isFolderSetKey(event), 1)
			} else {
				event.IsEdited = false
				r.sendAnswers(ctx, b, []*entities.Answer{sendMessage(event, message)})
			}
		}()
		return
	}
	var answers []*entities.Answer
	if message := r.process.MoveNoteEnd(ctx, event); message != "" {
		log.Debug("move note", logger.String("message", message))
		answers = append(answers, sendMessage(event, message))
	}
	log.Debug("select folder", logger.Int("folder_id", event.FolderID))
	notes, folderName := r.process.SelectFolder(ctx, event)
	answers = append(answers, sendMessage(event, messages.FolderEmoji))
	answers = append(answers, checkDefaultFolder(event, folderName))
	answers = collectNotes(answers, event, notes)

	go func() {
		r.deleteMessages(ctx, b, event)
		r.sendAnswers(ctx, b, answers)
	}()
}

func (r *router) doEmpty(ctx context.Context, b *bot.Bot, event *entities.Event) {
	if message := r.process.AddFolderEnd(ctx, event); message != "" {
		btns := r.process.Folders(ctx, event)
		event.IsEdited = true
		go func() {
			r.deleteMessages(ctx, b, event)
			r.sendAnswers(ctx, b, []*entities.Answer{sendFoldersList(event, btns, false)})
			r.process.SetInt(isFolderSetKey(event), 1)
		}()
		return
	}
	ap := r.process.Save(ctx, event)
	go func() {
		if ap.Message != "" {
			r.sendAnswers(ctx, b, []*entities.Answer{
				sendNote(event, event.NoteID, event.FolderID, true, ap),
			})
		}
	}()
}

func (r *router) doStart(ctx context.Context, b *bot.Bot, event *entities.Event) {
	message, btn := r.process.Start(ctx, event)
	event.Meta.MessageID = r.process.FolderMsgID(event.Meta.UserID)
	go func() {
		r.deleteMessages(ctx, b, event)
		r.deleteMessage(ctx, b, event)
		r.process.SetInt(isFolderSetKey(event), 0)
		r.sendAnswers(ctx, b, []*entities.Answer{sendFoldersButton(event, message, btn, true)})
	}()
}

func (r *router) doShowFolders(ctx context.Context, b *bot.Bot, event *entities.Event) {
	btns := r.process.Folders(ctx, event)
	isFolderSet := r.process.Int("isFolderSet:" + strconv.FormatInt(event.Meta.UserID, 10))
	go func() {
		r.deleteMessages(ctx, b, event)
		if isFolderSet == 0 {
			r.sendAnswers(ctx, b, []*entities.Answer{sendFoldersList(event, btns, false)})
			r.process.SetInt(isFolderSetKey(event), 1)
		}
	}()
}

func (r *router) doSaveTo(ctx context.Context, b *bot.Bot, event *entities.Event) {
	log := logger.L(ctx).With(logger.String("operation", "router.doSaveTo"))
	if message := r.process.SaveTo(ctx, event); message != "" {
		btns := r.process.Folders(ctx, event)
		event.Meta.MessageID = r.process.FolderMsgID(event.Meta.UserID)
		log.Debug(
			"save note to",
			logger.String("message", message),
			logger.Int("foldersID", event.Meta.MessageID),
		)
		go func() {
			r.deleteMessages(ctx, b, event)
			r.sendAnswers(ctx, b, []*entities.Answer{sendMessage(event, message)})
			event.IsEdited = true
			r.sendAnswers(ctx, b, []*entities.Answer{sendFoldersList(event, btns, false)})
			r.process.SetInt(isFolderSetKey(event), 1)
		}()
	}
}

func (r *router) doUnknown(ctx context.Context, b *bot.Bot, event *entities.Event) {
	r.sendAnswers(ctx, b, []*entities.Answer{
		sendMessage(event, messages.UnknownCommand),
	})
}

func (r *router) doMoveNote(ctx context.Context, b *bot.Bot, event *entities.Event) {
	event.NoteID, event.FolderID, _ = ParseButtonCallback(event.Text)
	message := r.process.MoveNoteStart(ctx, event)
	event.IsEdited = true
	go func() {
		r.deleteMessages(ctx, b, event)
		r.sendAnswers(ctx, b, []*entities.Answer{
			sendMessage(event, message),
		})
	}()
}

func (r *router) doDeleteNote(ctx context.Context, b *bot.Bot, event *entities.Event) {
	event.NoteID, event.FolderID, _ = ParseButtonCallback(event.Text)
	message := r.process.RemoveNote(ctx, event)
	go func() {
		r.deleteMessage(ctx, b, event)
		r.sendAnswers(ctx, b, []*entities.Answer{sendMessage(event, message)})
	}()
}

func (r *router) doCreateFolder(ctx context.Context, b *bot.Bot, event *entities.Event) {
	message := r.process.AddFolderStart(ctx, event)
	go r.sendAnswers(ctx, b, []*entities.Answer{sendMessage(event, message)})
}

func (r *router) doDeleteFolder(ctx context.Context, b *bot.Bot, event *entities.Event) {
	message := r.process.DeleteFolderStart(ctx, event)
	go r.sendAnswers(ctx, b, []*entities.Answer{sendMessage(event, message)})
}

func sendMessage(event *entities.Event, message string) *entities.Answer {
	return entities.NewAnswer(event, true, &entities.AnswerParams{Message: message})
}

func isFolderSetKey(event *entities.Event) string {
	return "isFolderSet:" + strconv.FormatInt(event.Meta.UserID, 10)
}

func (r *router) doInfo(ctx context.Context, b *bot.Bot, event *entities.Event) {
	videos := make([]*entities.Answer, 0, len(messages.InfoMap))
	for videoID, caption := range messages.InfoMap {
		ans := entities.Answer{UserID: event.Meta.UserID, DeleteAfter: true}
		ans.SendVideo = &bot.SendVideoParams{
			ChatID:  event.Meta.ChatID,
			Video:   &models.InputFileString{Data: videoID},
			Caption: caption,
		}
		videos = append(videos, &ans)
	}

	sort.Slice(videos, func(i, j int) bool {
		return videos[i].SendVideo.Caption < videos[j].SendVideo.Caption
	})

	go func() {
		r.deleteMessage(ctx, b, event)
		r.sendAnswers(ctx, b, videos)
	}()
}

func checkDefaultFolder(event *entities.Event, folderName string) *entities.Answer {
	if folderName != "default" {
		return sendFoldersButton(event, "📁"+folderName, buttons.Folders, true)
	}
	return sendFoldersButton(event, "📁"+buttons.DefaultFolderName, buttons.Folders, true)
}

func collectNotes(
	answers []*entities.Answer,
	event *entities.Event,
	notes map[int]*entities.AnswerParams,
) []*entities.Answer {
	if len(notes) == 0 {
		return append(answers, sendMessage(event, messages.NotesIsEmpty))
	}
	notesIDs := make([]int, 0, len(notes))
	for id := range notes {
		notesIDs = append(notesIDs, id)
	}
	sort.Ints(notesIDs)

	for _, noteID := range notesIDs {
		answers = append(answers, sendNote(
			event, noteID, event.FolderID, true, notes[noteID],
		))
	}

	return answers
}
