package router

import (
	"context"
	"fmt"

	"archive_bot/internal/entities"

	"archive_bot/pkg/logger"

	"github.com/go-telegram/bot"
)

func (r *router) sendAnswers(ctx context.Context, b *bot.Bot, answers []*entities.Answer) {
	log := logger.L(ctx).With(logger.String("operation", "router.sendAnswers"))

	for _, ans := range answers {
		ans := ans
		if ans.AnswerCallbackQuery != nil {
			go func() {
				if _, err := b.AnswerCallbackQuery(
					ctx, ans.AnswerCallbackQuery); err != nil {
					log.Error("AnswerCallbackQuery", logger.ErrAttr(err))
				}
			}()
		}

		if ans.SendMediaGroup != nil {
			messages, err := b.SendMediaGroup(ctx, ans.SendMediaGroup)
			checkErrAnswer("SendMediaGroup", log, err)
			for _, msg := range messages {
				r.checkIfMessageDeleteAfter(ans, msg.ID)
			}
		}
		if ans.SendPhoto != nil {
			msg, err := b.SendPhoto(ctx, ans.SendPhoto)
			checkErrAnswer("SendPhoto", log, err)
			r.checkIfMessageDeleteAfter(ans, msg.ID)
		}
		if ans.SendDocument != nil {
			msg, err := b.SendDocument(ctx, ans.SendDocument)
			checkErrAnswer("SendDocument", log, err)
			r.checkIfMessageDeleteAfter(ans, msg.ID)
		}
		if ans.SendVideo != nil {
			msg, err := b.SendVideo(ctx, ans.SendVideo)
			checkErrAnswer("SendVideo", log, err)
			r.checkIfMessageDeleteAfter(ans, msg.ID)
		}
		if ans.SendAudio != nil {
			msg, err := b.SendAudio(ctx, ans.SendAudio)
			checkErrAnswer("SendAudio", log, err)
			r.checkIfMessageDeleteAfter(ans, msg.ID)
		}
		if ans.SendAnimation != nil {
			msg, err := b.SendAnimation(ctx, ans.SendAnimation)
			checkErrAnswer("SendAnimation", log, err)
			r.checkIfMessageDeleteAfter(ans, msg.ID)
		}
		if ans.SendVoice != nil {
			msg, err := b.SendVoice(ctx, ans.SendVoice)
			checkErrAnswer("SendVoice", log, err)
			r.checkIfMessageDeleteAfter(ans, msg.ID)
		}
		if ans.SendMessage != nil {
			msg, err := b.SendMessage(ctx, ans.SendMessage)
			checkErrAnswer("SendMessage", log, err)
			r.checkIfMessageDeleteAfter(ans, msg.ID)
		}
		if ans.EditMessageText != nil {
			fmt.Println("EditMessageText", ans.EditMessageText)
			msg, err := b.EditMessageText(ctx, ans.EditMessageText)
			checkErrAnswer("EditMessageText", log, err)
			r.checkIfMessageDeleteAfter(ans, msg.ID)
		}
		if ans.EditMessageCaption != nil {
			msg, err := b.EditMessageCaption(ctx, ans.EditMessageCaption)
			checkErrAnswer("EditMessageCaption", log, err)
			r.checkIfMessageDeleteAfter(ans, msg.ID)
		}
		if ans.EditMessageMedia != nil {
			msg, err := b.EditMessageMedia(ctx, ans.EditMessageMedia)
			checkErrAnswer("EditMessageMedia", log, err)
			r.checkIfMessageDeleteAfter(ans, msg.ID)
		}
	}
}

// TODO: add context
func (r *router) checkIfMessageDeleteAfter(ans *entities.Answer, msgID int) {
	if ans.DeleteAfter {
		r.process.AddMessageID(ans.UserID, msgID)
	} else {
		r.log.Debug("checkIfMessageDeleteAfter", logger.Int("message ID", msgID))
		r.process.SetFolderMsgID(ans.UserID, msgID)
	}
}

func checkErrAnswer(method string, log *logger.Logger, err error) {
	if err != nil {
		log.Error(method, logger.ErrAttr(err))
	}
}
