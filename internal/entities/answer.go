package entities

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var captionLength int = 1024

type Answer struct {
	UserID              int64
	DeleteAfter         bool
	SendMessage         *bot.SendMessageParams
	SendPhoto           *bot.SendPhotoParams
	SendDocument        *bot.SendDocumentParams
	SendVideo           *bot.SendVideoParams
	SendAudio           *bot.SendAudioParams
	SendAnimation       *bot.SendAnimationParams
	SendVoice           *bot.SendVoiceParams
	SendMediaGroup      *bot.SendMediaGroupParams
	AnswerCallbackQuery *bot.AnswerCallbackQueryParams
	EditMessageText     *bot.EditMessageTextParams
	EditMessageCaption  *bot.EditMessageCaptionParams
	EditMessageMedia    *bot.EditMessageMediaParams
}

type AnswerParams struct {
	Type     Type
	Message  string
	FileIDs  []string
	Keyboard models.ReplyMarkup
}

func NewAnswer(event *Event, deleteAfter bool, ap *AnswerParams) *Answer {
	ans := &Answer{
		UserID:      event.Meta.UserID,
		DeleteAfter: deleteAfter,
	}

	if event.IsCallbackQuery {
		ans.AnswerCallbackQuery = &bot.AnswerCallbackQueryParams{
			CallbackQueryID: event.Meta.CallbackQueryID,
			ShowAlert:       false,
		}
	}

	if event.IsEdited {
		prepareEditedParams(ans, ap, event)
		return ans
	}

	switch ap.Type {
	case Photo:
		if len(ap.FileIDs) == 1 {
			ans.SendPhoto = &bot.SendPhotoParams{
				ChatID: event.Meta.ChatID,
				Photo:  &models.InputFileString{Data: ap.FileIDs[0]},
			}
			if len(ap.Message) < captionLength {
				ans.SendPhoto.Caption = ap.Message
				ans.SendPhoto.ReplyMarkup = ap.Keyboard
				return ans
			}
		} else {
			mediaArr := make([]models.InputMedia, 0, len(ap.FileIDs))
			for _, ph := range ap.FileIDs {
				mediaArr = append(mediaArr, &models.InputMediaPhoto{Media: ph})
			}
			ans.SendMediaGroup = &bot.SendMediaGroupParams{
				ChatID: event.Meta.ChatID,
				Media:  mediaArr,
			}
		}
	case Document:
		if len(ap.FileIDs) == 1 {
			ans.SendDocument = &bot.SendDocumentParams{
				ChatID:   event.Meta.ChatID,
				Document: &models.InputFileString{Data: ap.FileIDs[0]},
			}
			if len(ap.Message) < captionLength {
				ans.SendDocument.Caption = ap.Message
				ans.SendDocument.ReplyMarkup = ap.Keyboard
				return ans
			}
		} else {
			mediaArr := make([]models.InputMedia, 0, len(ap.FileIDs))
			for _, doc := range ap.FileIDs {
				mediaArr = append(mediaArr, &models.InputMediaDocument{Media: doc})
			}
			ans.SendMediaGroup = &bot.SendMediaGroupParams{
				ChatID: event.Meta.ChatID,
				Media:  mediaArr,
			}
		}
	case Video:
		if len(ap.FileIDs) == 1 {
			ans.SendVideo = &bot.SendVideoParams{
				ChatID: event.Meta.ChatID,
				Video:  &models.InputFileString{Data: ap.FileIDs[0]},
			}
			if len(ap.Message) < captionLength {
				ans.SendVideo.Caption = ap.Message
				ans.SendVideo.ReplyMarkup = ap.Keyboard
				return ans
			}
		} else {
			mediaArr := make([]models.InputMedia, 0, len(ap.FileIDs))
			for _, vid := range ap.FileIDs {
				mediaArr = append(mediaArr, &models.InputMediaVideo{Media: vid})
			}
			ans.SendMediaGroup = &bot.SendMediaGroupParams{
				ChatID: event.Meta.ChatID,
				Media:  mediaArr,
			}
		}
	case Audio:
		if len(ap.FileIDs) == 1 {
			ans.SendAudio = &bot.SendAudioParams{
				ChatID: event.Meta.ChatID,
				Audio:  &models.InputFileString{Data: ap.FileIDs[0]},
			}
			if len(ap.Message) < captionLength {
				ans.SendAudio.Caption = ap.Message
				ans.SendAudio.ReplyMarkup = ap.Keyboard
				return ans
			}
		} else {
			mediaArr := make([]models.InputMedia, 0, len(ap.FileIDs))
			for _, a := range ap.FileIDs {
				mediaArr = append(mediaArr, &models.InputMediaAudio{Media: a})
			}
			ans.SendMediaGroup = &bot.SendMediaGroupParams{
				ChatID: event.Meta.ChatID,
				Media:  mediaArr,
			}
		}
	case Animation:
		if len(ap.FileIDs) == 1 {
			ans.SendAnimation = &bot.SendAnimationParams{
				ChatID:    event.Meta.ChatID,
				Animation: &models.InputFileString{Data: ap.FileIDs[0]},
			}
			if len(ap.Message) < captionLength {
				ans.SendAnimation.Caption = ap.Message
				ans.SendAnimation.ReplyMarkup = ap.Keyboard
				return ans
			}
		} else {
			mediaArr := make([]models.InputMedia, 0, len(ap.FileIDs))
			for _, ani := range ap.FileIDs {
				mediaArr = append(mediaArr, &models.InputMediaAnimation{Media: ani})
			}
			ans.SendMediaGroup = &bot.SendMediaGroupParams{
				ChatID: event.Meta.ChatID,
				Media:  mediaArr,
			}
		}
	case Voice:
		if len(ap.FileIDs) == 1 {
			ans.SendVoice = &bot.SendVoiceParams{
				ChatID: event.Meta.ChatID,
				Voice:  &models.InputFileString{Data: ap.FileIDs[0]},
			}
			if len(ap.Message) < captionLength {
				ans.SendVoice.Caption = ap.Message
				ans.SendVoice.ReplyMarkup = ap.Keyboard
				return ans
			}
		}
	}

	ans.SendMessage = &bot.SendMessageParams{
		ChatID:      event.Meta.ChatID,
		Text:        ap.Message,
		ReplyMarkup: ap.Keyboard,
	}

	return ans
}

func prepareEditedParams(ans *Answer, ap *AnswerParams, event *Event) {
	ans.EditMessageText = &bot.EditMessageTextParams{
		ChatID:      event.Meta.ChatID,
		MessageID:   event.Meta.MessageID,
		Text:        ap.Message,
		ReplyMarkup: ap.Keyboard,
	}
}
