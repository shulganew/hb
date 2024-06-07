package bot

import (
	"context"
	"fmt"
	"log"
	"reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shulganew/hb.git/internal/config"
)

func BotHandler(ctx context.Context, conf config.Config, componentsErrs chan error, botDone chan struct{}) {
	// Create new bot.
	bot, err := tgbotapi.NewBotAPI(conf.Bot)
	if err != nil {
		componentsErrs <- fmt.Errorf("listen and server has failed: %w", err)
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
			if update.Message != nil { // If we got a message
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
					switch update.Message.Text {
					case "/start":
						// Intro.
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, I'm a happy birthday bot, I remind you to greet your friends!")
						bot.Send(msg)

					case "/list":
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello")
						msg.ReplyToMessageID = update.Message.MessageID
						fmt.Println("----")
						fmt.Printf("%#v", update.Message.From)
						fmt.Println("----")
						fmt.Printf("%#v", update.MyChatMember)
						fmt.Println("----")
						fmt.Printf("%#v", update.SentFrom())
						fmt.Println("----")
						fmt.Printf("%#v", update.ChatMember)
						fmt.Println("----")

						bot.Send(msg)
					case "/login":
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
					// Respond to the callback query, telling Telegram to show the user
					// a message with the data received.
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
					if _, err := bot.Request(callback); err != nil {
						componentsErrs <- fmt.Errorf("listen and server has failed: %w", err)
						return
					}

					// And finally, send a message containing the data received.
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
					if _, err := bot.Send(msg); err != nil {
						componentsErrs <- fmt.Errorf("listen and server has failed: %w", err)
						return
					}
				}
			}
		}
	}
}
