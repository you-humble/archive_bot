package app

import (
	"context"
	"archive_bot/pkg/closer"
	"archive_bot/pkg/er"
	"archive_bot/pkg/logger"
	"net/http"
	"runtime"

	"github.com/go-telegram/bot"
)

type app struct {
	dp      *dependencyProvider
	bot     *bot.Bot
	workers int
}

func New(ctx context.Context, workers int) *app {
	if workers == 0 {
		workers = runtime.NumCPU()
	}

	a := &app{workers: workers}
	a.initBot(ctx)

	return a
}

func (a *app) initBot(ctx context.Context) {
	a.dp = newDependencyProvider()

	ctx = logger.ContextWithLogger(ctx, a.dp.Logger())

	opts := []bot.Option{
		bot.WithWorkers(a.workers),
		bot.WithDefaultHandler(
			a.dp.Router(ctx).RouteMessage,
		),
		bot.WithCallbackQueryDataHandler(
			"btn_", bot.MatchTypePrefix,
			a.dp.Router(ctx).RouteCallbackQuery,
		),
		bot.WithMessageTextHandler(
			"adm", bot.MatchTypeContains,
			a.dp.Router(ctx).RouteAdminMessage,
		),
		bot.WithCallbackQueryDataHandler(
			"adm", bot.MatchTypeContains,
			a.dp.Router(ctx).RouteAdminCallback),
	}

	b, err := bot.New(a.dp.Config().Bot.Token, opts...)
	if err != nil {
		panic(er.New("failed to initialize the bot", "app.initBot", err))
	}

	a.bot = b
}

func (a *app) Run(ctx context.Context) {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	a.selectConnection(ctx)

	go func() {
		if err := http.ListenAndServe(
			":"+a.dp.Config().Bot.Port,
			a.bot.WebhookHandler(),
		); err != nil {
			a.dp.Logger().Error("Bot stopped due error", logger.ErrAttr(err))
			return
		}
	}()

	a.dp.Logger().Info(
		"Work",
		logger.String("port", a.dp.Config().Bot.Port),
		logger.Int("workers", a.workers),
	)
	a.dp.Logger().Info("✓ Bot started!")
	if a.dp.Config().IsWebhook == 1 {
		a.startWebhook(ctx)
	} else {
		a.bot.Start(ctx)
	}
	a.dp.Logger().Info("𐄂 Bot stopped!")
}

func (a *app) startWebhook(ctx context.Context) {
	a.bot.StartWebhook(ctx)
}

func (a *app) selectConnection(ctx context.Context) {
	if a.dp.Config().IsWebhook == 1 {
		if _, err := a.bot.SetWebhook(ctx, &bot.SetWebhookParams{
			URL:                a.dp.Config().Bot.WebhookURL,
			DropPendingUpdates: true,
		}); err != nil {
			a.dp.Logger().Error("SetWebhook error:", logger.ErrAttr(err))
			panic("SetWebhook error:" + err.Error())
		}
		a.dp.Logger().Info("webhook is set", logger.String("url", a.dp.Config().Bot.WebhookURL))
	} else {
		if _, err := a.bot.DeleteWebhook(ctx, &bot.DeleteWebhookParams{
			DropPendingUpdates: true,
		}); err != nil {
			a.dp.Logger().Error("DeleteWebhook error:", logger.ErrAttr(err))
			panic("DeleteWebhook error:" + err.Error())
		}
		a.dp.Logger().Info("polling is set")
	}
}
