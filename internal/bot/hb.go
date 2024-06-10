package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shulganew/hb.git/internal/config"
	"github.com/shulganew/hb.git/internal/entities"
	"github.com/shulganew/hb.git/internal/services"
	"go.uber.org/zap"
)

func BotHandler(ctx context.Context, conf config.Config, b *tgbotapi.BotAPI, bs *services.Bot, componentsErrs chan error, botDone chan struct{}) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)

	for {
		select {
		// Contex done, exit,
		case <-ctx.Done():
			// Exit on errors.
			close(botDone)
			return
		// Read bot updates.
		case update := <-updates:
			if update.Message != nil {
				// Check is loged in user.
				if !bs.IsLogedIn(update.Message.From.UserName) {
					zap.S().Debugln("User not loged in: ", update.Message.From.UserName)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для начала работы выполните регистрацю.")
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonLoginURL("Нажми, чтобы зарегестрироваться", tgbotapi.LoginURL{URL: config.Domain}),
						))
					b.Send(msg)
					continue
				}

				switch update.Message.Text {
				case "/start":
					// Intro.
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Этот бот напомнит поздравить с днем рождения!")
					b.Send(msg)

				// List all available users.
				case "/list":
					all := bs.ListAll()
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, all)
					b.Send(msg)

				// Subscribe to users happy birthday.
				case "/sub":
					// Get all usres available for subscribtion.
					allSub := bs.ListAvalibleSub(update.Message.From.UserName)

					// Constract answer with inline buttons.
					var buttons []tgbotapi.InlineKeyboardButton
					if len(allSub) == 0 {
						b.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Вы подписаны на всех одступных вам пользователей."))
					}

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберете пользователя для подписки на др:")
					for _, user := range allSub {
						// Prepare action
						action := entities.Action{Type: entities.SUB, ChatID: update.Message.Chat.ID, Subscriber: update.Message.From.UserName, Subscribed: user.Name}
						eAction := services.EndcodeAction(action)

						zap.S().Debugln(eAction)
						buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(user.Name, eAction))
					}
					var rows [][]tgbotapi.InlineKeyboardButton
					for _, button := range buttons {
						rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
					}
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
					b.Send(msg)

				case "/unsub":
					// Get all subscribtions.
					allSub := bs.ListCurrentSub(update.Message.From.UserName)

					// Constract answer with inline buttons.
					var buttons []tgbotapi.InlineKeyboardButton
					if len(allSub) == 0 {
						b.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Вы не подписаны на пользователей."))
					}

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отказаться от подписки:")
					for _, user := range allSub {
						// Prepare action
						action := entities.Action{Type: entities.UNSUB, ChatID: update.Message.Chat.ID, Subscriber: update.Message.From.UserName, Subscribed: user.Name}
						eAction := services.EndcodeAction(action)

						buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(user.Name, eAction))
					}
					var rows [][]tgbotapi.InlineKeyboardButton
					for _, button := range buttons {
						rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
					}
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
					b.Send(msg)
				default:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi all!")
					b.Send(msg)
				}
			} else if update.CallbackQuery != nil {
				// Get answer from Bot service.
				zap.S().Debugln("Callback: ", update.CallbackQuery.Data)
				ansver := bs.Resiver(update.CallbackQuery.Data)

				// And finally, send a message containing the data received.
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ansver)
				if _, err := b.Send(msg); err != nil {
					componentsErrs <- fmt.Errorf("bot failed: %w", err)
					return
				}
			}
		}
	}
}
