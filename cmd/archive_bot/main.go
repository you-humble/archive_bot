package main

import (
	"context"
	"os"
	"os/signal"
	
	"archive_bot/internal/app"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	telegramBot := app.New(ctx, 0)
	telegramBot.Run(ctx)
}
