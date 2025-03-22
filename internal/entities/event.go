package entities

import (
	"context"
	"archive_bot/pkg/logger"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot/models"
)

type Event struct {
	Type            Type
	IsCallbackQuery bool
	IsEdited        bool
	Text            string
	FileID          string
	MediaGroupID    string
	NoteID          int
	FolderID        int
	Meta            Meta
}

type Meta struct {
	UserID          int64
	ChatID          int64
	MessageID       int
	UserName        string
	CallbackQueryID string
	Date            time.Time
	Media           string
}

func NewEvent(ctx context.Context, update *models.Update) *Event {
	log := logger.L(ctx).With(logger.String("operation", "entities.NewEvent"))

	event := &Event{}
	eventType := fetchType(update, event)
	event.Type = eventType

	if update.Message != nil {
		log.Debug(
			"update data",
			logger.String("text", update.Message.Text),
			logger.String("Caption", update.Message.Caption),
		)
	}

	if event.IsCallbackQuery {
		event.Text = update.CallbackQuery.Data
		event.Meta = Meta{
			UserID:          update.CallbackQuery.From.ID,
			ChatID:          update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:       update.CallbackQuery.Message.Message.ID,
			UserName:        update.CallbackQuery.From.Username,
			CallbackQueryID: update.CallbackQuery.ID,
			Date:            time.Unix(int64(update.CallbackQuery.Message.Message.Date), 0),
		}
		return event
	}

	switch eventType {
	case Message:
		event.Text = checkForwardOrigin(update, event)
	case Photo:
		event.Text = checkForwardOrigin(update, event)
		event.MediaGroupID = update.Message.MediaGroupID
		photo := update.Message.Photo[len(update.Message.Photo)-1]
		event.FileID = photo.FileID
	case Document:
		event.Text = checkForwardOrigin(update, event)
		event.MediaGroupID = update.Message.MediaGroupID
		event.FileID = update.Message.Document.FileID
	case Video:
		event.Text = checkForwardOrigin(update, event)
		event.MediaGroupID = update.Message.MediaGroupID
		event.FileID = update.Message.Video.FileID
	case Audio:
		event.MediaGroupID = update.Message.MediaGroupID
		event.FileID = update.Message.Audio.FileID
	case Animation:
		event.MediaGroupID = update.Message.MediaGroupID
		event.FileID = update.Message.Animation.FileID
	case Voice:
		event.Text = checkForwardOrigin(update, event)
		event.FileID = update.Message.Voice.FileID
	}
	event.Meta = Meta{
		UserID:    update.Message.From.ID,
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
		UserName:  update.Message.From.Username,
		Date:      time.Unix(int64(update.Message.Date), 0),
	}

	log.Debug("new event data", logger.String("event", event.String()))
	return event
}

func checkForwardOrigin(update *models.Update, event *Event) string {
	if event.Type == Message {
		if update.Message.ForwardOrigin != nil {
			return setSource(update, update.Message.Text)
		} else {
			return update.Message.Text
		}
	}
	if update.Message.ForwardOrigin != nil {
		return setSource(update, update.Message.Caption)
	} else {
		return update.Message.Caption
	}
}

const source string = "Источник: @"

func setSource(update *models.Update, text string) string {
	b := &strings.Builder{}
	b.WriteString(source)
	messageOrigin := update.Message.ForwardOrigin
	switch {
	case messageOrigin.MessageOriginChannel != nil:
		if messageOrigin.MessageOriginChannel.Chat.Username != "" {
			b.WriteString(messageOrigin.MessageOriginChannel.Chat.Username)
		}
	case messageOrigin.MessageOriginChat != nil:
		if messageOrigin.MessageOriginChat.SenderChat.Username != "" {
			b.WriteString(messageOrigin.MessageOriginChat.SenderChat.Username)
		}
	case messageOrigin.MessageOriginHiddenUser != nil:
		if messageOrigin.MessageOriginHiddenUser.SenderUserName != "" {
			b.WriteString(messageOrigin.MessageOriginHiddenUser.SenderUserName)
		}
	case messageOrigin.MessageOriginChannel != nil:
		if messageOrigin.MessageOriginUser.SenderUser.Username != "" {
			b.WriteString(messageOrigin.MessageOriginUser.SenderUser.Username)
		}
	}

	if b.String() == source && text == "" {
		return ""
	}

	b.WriteString("\n\n")
	b.WriteString(text)
	return b.String()
}

func fetchType(update *models.Update, event *Event) Type {
	if update.CallbackQuery != nil {
		event.IsCallbackQuery = true
		if update.CallbackQuery.Message.Message != nil {
			if len(update.CallbackQuery.Message.Message.Photo) != 0 {
				return Photo
			}
			return Message
		}
	}
	switch {
	case len(update.Message.Photo) != 0:
		return Photo
	case update.Message.Document != nil:
		return Document
	case update.Message.Video != nil:
		return Video
	case update.Message.Audio != nil:
		return Audio
	case update.Message.Animation != nil:
		return Animation
	case update.Message.Voice != nil:
		return Voice
	case update.Message.Text != "":
		return Message
	}

	return Unknown
}

func (e *Event) String() string {
	b := &strings.Builder{}

	b.WriteString("Event{Type: ")
	b.WriteString(e.Type.String())
	b.WriteString(", IsCallbackQuery: ")
	b.WriteString(strconv.FormatBool(e.IsCallbackQuery))
	b.WriteString(", IsEdited: ")
	b.WriteString(strconv.FormatBool(e.IsEdited))
	b.WriteString(", Text: ")
	b.WriteString(e.Text)
	b.WriteString(", FileID: ")
	b.WriteString(e.FileID)
	b.WriteString(", MediaGroupID: ")
	b.WriteString(e.MediaGroupID)
	b.WriteString(", NoteID: ")
	b.WriteString(strconv.Itoa(e.NoteID))
	b.WriteString(", FolderID: ")
	b.WriteString(strconv.Itoa(e.FolderID))
	b.WriteString(", Meta: ")
	b.WriteString(e.Meta.String())
	b.WriteRune('}')

	return b.String()
}

func (m Meta) String() string {
	b := &strings.Builder{}

	b.WriteString("Meta{UserID: ")
	b.WriteString(strconv.FormatInt(m.UserID, 10))
	b.WriteString(", ChatID: ")
	b.WriteString(strconv.FormatInt(m.ChatID, 10))
	b.WriteString(", MessageID: ")
	b.WriteString(strconv.Itoa(m.MessageID))
	b.WriteString(", UserName: ")
	b.WriteString(m.UserName)
	b.WriteString(", CallbackQueryID: ")
	b.WriteString(m.CallbackQueryID)
	b.WriteString(", Date: ")
	b.WriteString(m.Date.String())
	b.WriteString(", Media: ")
	b.WriteString(m.Media)
	b.WriteRune('}')

	return b.String()
}
