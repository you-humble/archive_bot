package router

import (
	"context"
	"archive_bot/internal/entities"
	"archive_bot/pkg/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	admin      string = "/adm"
	usersCount string = "adm_count_users"
)

func (r *router) RouteAdminMessage(ctx context.Context, b *bot.Bot, update *models.Update) {
	log := r.log.With(logger.String("operation", "router.RouteAdmin"))

	if update.Message.From.ID != r.adminID {
		log.Warn("not admin wanted to use", logger.Int64("user_ID", update.Message.From.ID))
		return
	}

	event := entities.NewEvent(ctx, update)
	r.process.AddMessageID(event.Meta.UserID, event.Meta.MessageID)

	command, text := commandAndText(event.Text)
	if text != "" {
		event.Text = text
	}

	log.Debug("switch admin message", logger.String("command", command))
	switch command {
	case admin:
		go r.adminPanel(ctx, b, event)
	default:
		go r.doUnknown(ctx, b, event)
	}
}

func (r *router) RouteAdminCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	log := r.log.With(logger.String("operation", "router.RouteAdmin"))

	event := entities.NewEvent(ctx, update)
	r.process.AddMessageID(event.Meta.UserID, event.Meta.MessageID)

	if event.Meta.UserID != r.adminID {
		log.Warn("not admin wanted to use", logger.Int64("user_ID", update.Message.From.ID))
		return
	}

	log.Debug("switch admin callback", logger.String("command", event.Text))
	switch event.Text {
	case usersCount:
		go r.doCountUsers(ctx, b, event)

	}
}

func (r *router) adminPanel(ctx context.Context, b *bot.Bot, event *entities.Event) {
	btns := make([][]models.InlineKeyboardButton, 0, 1)
	btns = append(btns, []models.InlineKeyboardButton{
		{CallbackData: usersCount, Text: "Count of users"},
	})

	panel := entities.NewAnswer(event, true, &entities.AnswerParams{
		Message:  "Welcome to admin panel!",
		Keyboard: &models.InlineKeyboardMarkup{InlineKeyboard: btns},
	})

	go r.sendAnswers(ctx, b, []*entities.Answer{panel})
}

func (r *router) doCountUsers(ctx context.Context, b *bot.Bot, event *entities.Event) {
	message := r.process.CountUsers(ctx)
	go r.sendAnswers(ctx, b, []*entities.Answer{
		sendMessage(event, message),
	})
}
