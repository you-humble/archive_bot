package router

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"archive_bot/internal/const/buttons"
	"archive_bot/internal/const/messages"
	"archive_bot/internal/entities"

	"archive_bot/pkg/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	start             string = "/start"
	info              string = "/info"
	folders           string = "/folders"
	moveLastNote      string = "/move_note"
	moveLastNoteAlias string = "!"
)

type Processor interface {
	InitUser(ctx context.Context, event *entities.Event)
	CountUsers(ctx context.Context) string
	SetFolderMsgID(userID int64, messageID int)
	FolderMsgID(userID int64) int
	SetInt(key string, num int)
	Int(key string) int
	AddMessageID(userID int64, messageID int)
	MessageIDs(userID int64) []int

	Start(ctx context.Context, event *entities.Event) (string, string)
	Folders(ctx context.Context, event *entities.Event) map[string]string
	Save(ctx context.Context, event *entities.Event) *entities.AnswerParams
	SaveTo(ctx context.Context, event *entities.Event) string

	SelectFolder(ctx context.Context, event *entities.Event) (map[int]*entities.AnswerParams, string)
	AddFolderStart(ctx context.Context, event *entities.Event) string
	AddFolderEnd(ctx context.Context, event *entities.Event) string
	DeleteFolderStart(ctx context.Context, event *entities.Event) string
	DeleteFolderEnd(ctx context.Context, event *entities.Event) string

	MoveNoteStart(ctx context.Context, event *entities.Event) string
	MoveNoteEnd(ctx context.Context, event *entities.Event) string
	RemoveNote(ctx context.Context, event *entities.Event) string
}

type router struct {
	log     *logger.Logger
	adminID int64
	process Processor
}

func New(log *logger.Logger, adminID int64, processor Processor) *router {
	return &router{log: log, process: processor, adminID: adminID}
}

func (r *router) RouteCallbackQuery(ctx context.Context, b *bot.Bot, update *models.Update) {
	log := r.log.With(logger.String("operation", "router.RouteCallbackQuery"))

	event := entities.NewEvent(ctx, update)
	// r.process.AddMessageID(event.Meta.UserID, event.Meta.MessageID)
	r.process.InitUser(ctx, event)

	log.Debug("switch CallbackQuery", logger.String("command", event.Text))
	switch {
	case event.Text == buttons.CreateFolder:
		go r.doCreateFolder(ctx, b, event)
	case event.Text == buttons.DeleteFolder:
		go r.doDeleteFolder(ctx, b, event)
	case strings.HasPrefix(event.Text, buttons.DeleteNote):
		go r.doDeleteNote(ctx, b, event)
	case strings.HasPrefix(event.Text, buttons.MoveNote):
		go r.doMoveNote(ctx, b, event)
	default:
		go r.doDefaultCallback(ctx, b, event)
	}
}

func (r *router) RouteMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	log := r.log.With(logger.String("operation", "router.RouteMessage"))

	event := entities.NewEvent(ctx, update)
	r.process.AddMessageID(event.Meta.UserID, event.Meta.MessageID)
	r.process.InitUser(ctx, event)

	command, text := commandAndText(event.Text)
	if text != "" {
		event.Text = text
	}

	if event.Type == entities.Unknown {
		go r.doUnknown(ctx, b, event)
		return
	}

	log.Debug("switch Message", logger.String("command", command))
	switch command {
	case "":
		go r.doEmpty(ctx, b, event)
	case start:
		go r.doStart(ctx, b, event)
	case info:
		go r.doInfo(ctx, b, event)
	case folders:
		go r.doShowFolders(ctx, b, event)
	case moveLastNote:
		go r.doSaveTo(ctx, b, event)
	default:
		go r.doUnknown(ctx, b, event)
	}
}

func sendFoldersList(
	event *entities.Event,
	buttonsMap map[string]string,
	deleteAfter bool,
) *entities.Answer {
	btns := make([][]models.InlineKeyboardButton, 0, len(buttonsMap)+1)

	for key, val := range buttonsMap {
		btns = append(btns, []models.InlineKeyboardButton{
			{CallbackData: key, Text: val},
		})
	}

	sort.Slice(btns, func(i, j int) bool {
		return btns[i][0].CallbackData < btns[j][0].CallbackData
	})

	btns = append(btns, []models.InlineKeyboardButton{
		{
			CallbackData: buttons.CreateFolder,
			Text:         buttons.MenuOptions[buttons.CreateFolder],
		},
		{
			CallbackData: buttons.DeleteFolder,
			Text:         buttons.MenuOptions[buttons.DeleteFolder],
		},
	})

	return entities.NewAnswer(
		event,
		deleteAfter,
		&entities.AnswerParams{
			Message:  messages.FoldersCaption,
			Keyboard: &models.InlineKeyboardMarkup{InlineKeyboard: btns},
		})
}

func sendNote(
	event *entities.Event,
	noteID int,
	folderID int,
	deleteAfter bool,
	ap *entities.AnswerParams,
) *entities.Answer {
	btns := make([][]models.InlineKeyboardButton, 0, 1)

	buttonsRow := make([]models.InlineKeyboardButton, 0, len(buttons.CatalogueOptions))

	noteAndFolder := strconv.Itoa(noteID) + buttons.Delimiter + strconv.Itoa(folderID)

	switch ap.Type {
	case entities.Photo:
		if len(ap.FileIDs) == 1 {
			noteAndFolder += buttons.Delimiter + buttons.WithPhoto
		}
	}

	for key, val := range buttons.CatalogueOptions {
		buttonsRow = append(buttonsRow, models.InlineKeyboardButton{
			CallbackData: key + buttons.Delimiter + noteAndFolder,
			Text:         val,
		})
	}

	sort.Slice(buttonsRow, func(i, j int) bool {
		return buttonsRow[i].CallbackData < buttonsRow[j].CallbackData
	})

	btns = append(btns, buttonsRow)
	ap.Keyboard = &models.InlineKeyboardMarkup{
		InlineKeyboard: btns,
	}

	return entities.NewAnswer(event, deleteAfter, ap)
}

func sendFoldersButton(
	event *entities.Event,
	description string,
	button string,
	deleteAfter bool,
) *entities.Answer {
	kb := &models.ReplyKeyboardMarkup{
		Keyboard:       [][]models.KeyboardButton{{{Text: button}}},
		ResizeKeyboard: true,
	}

	return entities.NewAnswer(
		event,
		deleteAfter,
		&entities.AnswerParams{
			Message:  description,
			Keyboard: kb,
		})
}

func commandAndText(plainMessage string) (string, string) {
	if plainMessage == "" {
		return "", ""
	}

	if strings.HasPrefix(plainMessage, moveLastNoteAlias) {
		return moveLastNote, strings.TrimLeft(plainMessage, moveLastNoteAlias)
	}

	if plainMessage == buttons.Folders {
		return folders, ""
	}

	if !strings.HasPrefix(plainMessage, "/") {
		return "", plainMessage
	}

	count := 0
	for _, ch := range plainMessage {
		if ch != '/' {
			break
		}
		count++
	}
	if count > 1 {
		return "", plainMessage
	}

	strSlice := strings.SplitN(plainMessage, " ", 2)

	command := strSlice[0]
	if len(strSlice) == 1 {
		return command, ""
	}

	text := strings.TrimSpace(strSlice[1])

	return command, text
}

func ParseButtonCallback(command string) (int, int, string) {
	if command == "" {
		return 0, 0, ""
	}
	sl := strings.Split(command, buttons.Delimiter)
	if len(sl) < 3 {
		return 0, 0, ""
	}

	noteID, err := strconv.Atoi(sl[1])
	if err != nil {
		return 0, 0, ""
	}
	folderID, err := strconv.Atoi(sl[2])
	if err != nil {
		return 0, 0, ""
	}

	if len(sl) == 4 {
		withMedia := sl[3]
		return noteID, folderID, withMedia
	}

	return noteID, folderID, ""
}

func (r *router) deleteMessages(ctx context.Context, b *bot.Bot, event *entities.Event) {
	msgIDs := r.process.MessageIDs(event.Meta.UserID)
	if len(msgIDs) > 0 {
		_, err := b.DeleteMessages(ctx, &bot.DeleteMessagesParams{
			ChatID:     event.Meta.ChatID,
			MessageIDs: msgIDs,
		})
		if err != nil {
			r.log.Error("", logger.ErrAttr(err))
		}
	}
}

func (r *router) deleteMessage(ctx context.Context, b *bot.Bot, event *entities.Event) {
	if event.Meta.MessageID != 0 {
		_, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
			ChatID:    event.Meta.ChatID,
			MessageID: event.Meta.MessageID,
		})
		if err != nil {
			r.log.Error("", logger.ErrAttr(err))
		}
	}
}
