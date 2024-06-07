package app

import (
	"context"

	"github.com/shulganew/hb.git/internal/bot"
	"github.com/shulganew/hb.git/internal/config"
)

// Telegram bot api handler.
func StartBot(ctx context.Context, conf config.Config, componentsErrs chan error) (botDone chan struct{}) {

	// Graceful shutdown.
	botDone = make(chan struct{})

	// Start bot handling.
	go bot.BotHandler(ctx, conf, componentsErrs, botDone)

	return

}
