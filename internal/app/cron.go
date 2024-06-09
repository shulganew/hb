package app

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"github.com/shulganew/hb.git/internal/config"
	"github.com/shulganew/hb.git/internal/services"
	"go.uber.org/zap"
)

func InitCron(ctx context.Context, b *tgbotapi.BotAPI, bs *services.Bot) (c *cron.Cron) {
	c = cron.New()
	c.AddFunc(config.Cron, func() {
		users := bs.CheckHB(ctx)
		for _, user := range users {
			zap.S().Infof("user %#v \n", user)
			// Make notification for all subscribers in telegram chats.
			chats := bs.GetNOtifyChats(ctx, user)
			for _, chatID := range chats {
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Поздравляю %s с др!", user.Name))
				b.Send(msg)
			}
		}
	})
	c.Start()
	return
}
