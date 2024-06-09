package bot

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shulganew/hb.git/internal/config"
	"github.com/shulganew/hb.git/internal/entities"
	"github.com/shulganew/hb.git/internal/services"
	"go.uber.org/zap"
)

func BotHandler(ctx context.Context, conf config.Config, bs *services.Bot, componentsErrs chan error, botDone chan struct{}) {
	// Create new bot.
	bot, err := tgbotapi.NewBotAPI(conf.Bot)
	if err != nil {
		componentsErrs <- fmt.Errorf("bot failed: %w", err)
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

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
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для начала работы выполните регистрацю.")
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonLoginURL("Нажми, чтобы зарегестрироваться", tgbotapi.LoginURL{URL: "https://learn.iskratechno.ru"}),
						))
					bot.Send(msg)
					continue
				}

				switch update.Message.Text {
				case "/start":
					// Intro.
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Этот бот напомнит поздравить с днем рождения!")
					bot.Send(msg)

				// List all available users.
				case "/list":
					all := bs.ListAll()
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, all)
					bot.Send(msg)

				// Subscribe to users happy birthday.
				case "/sub":
					// Get all usres available for subscribtion.
					allSub := bs.ListAvalibleSub(update.Message.From.UserName)

					// Constract answer with inline buttons.
					var buttons []tgbotapi.InlineKeyboardButton
					if len(allSub) == 0 {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Вы подписаны на всех одступных вам пользователей."))
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
					bot.Send(msg)

				case "/unsub":
					// Get all subscribtions.
					allSub := bs.ListCurrentSub(update.Message.From.UserName)

					// Constract answer with inline buttons.
					var buttons []tgbotapi.InlineKeyboardButton
					if len(allSub) == 0 {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Вы не подписаны на пользователей."))
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
					bot.Send(msg)

				case "/info":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello")
					msg.ReplyToMessageID = update.Message.MessageID
					zap.S().Infoln("----")
					zap.S().Infof("%#v", update.Message.From)
					zap.S().Infoln("----")
					zap.S().Infof("%#v", update.MyChatMember)
					zap.S().Infoln("----")
					zap.S().Infof("%#v", update.SentFrom())
					zap.S().Infoln("----")
					zap.S().Infof("%#v", update.ChatMember)
					zap.S().Infoln("----")

					bot.Send(msg)
				case "/login":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "d")
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonLoginURL("rrr", tgbotapi.LoginURL{URL: "https://learn.iskratechno.ru"}),
						))
					bot.Send(msg)
				case "/click":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "data")
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("data", "mydata"),
						))
					bot.Send(msg)
				case "/test":
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "d")
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonLoginURL("rrr", tgbotapi.LoginURL{URL: "https://learn.iskratechno.ru"}),
						))
					bot.Send(msg)
				case "/invite":
					invite, err := bot.GetInviteLink(tgbotapi.ChatInviteLinkConfig{ChatConfig: update.FromChat().ChatConfig()})
					if err != nil {
						log.Println(err)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
						bot.Send(msg)
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, invite)
					bot.Send(msg)
				case "/contact":
					conf := tgbotapi.NewContact(update.Message.Chat.ID, "+9996621111", "Igor")
					bot.Send(conf)
				default:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi all!")
					bot.Send(msg)
				}
			} else if update.CallbackQuery != nil {
				// Get answer from Bot service.
				zap.S().Infoln("Callback: ", update.CallbackQuery.Data)
				ansver := bs.Resiver(update.CallbackQuery.Data)

				// And finally, send a message containing the data received.
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, ansver)
				if _, err := bot.Send(msg); err != nil {
					componentsErrs <- fmt.Errorf("bot failed: %w", err)
					return
				}
			}
		}
	}
}
