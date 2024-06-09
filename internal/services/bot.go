package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shulganew/hb.git/internal/config"
	"github.com/shulganew/hb.git/internal/entities"
	"go.uber.org/zap"
)

// User creation, registration
type Bot struct {
	stor    Happyer
	conf    config.Config
	logedIn map[string]struct{} // Autentificated telegram users.
}

func NewBot(ctx context.Context, stor Happyer, conf config.Config) *Bot {
	ln := make(map[string]struct{}, 0)

	// TODO remove!
	ln["shulgaigor"] = struct{}{}

	bs := &Bot{stor: stor, conf: conf, logedIn: ln}
	return bs
}

// Check if user loged in.
func (b Bot) IsLogedIn(tg string) (isLogedIn bool) {
	_, ok := b.logedIn[tg]
	if ok {
		return true
	}
	return
}

// List all existed users.
func (b Bot) ListAll() string {
	users, err := b.stor.ListAll(context.Background())
	if err != nil {
		zap.S().Errorln(err)
		return ""
	}

	var ans strings.Builder
	for _, user := range users {
		ans.WriteString("- ")
		ans.WriteString(user.Name)
		ans.WriteString(" ")
		ans.WriteString(user.Hb.Format("02-01-2006"))
		ans.WriteString("\n")
	}
	return ans.String()
}

// List all existed users for subscribtion for user..
func (b Bot) ListAvalibleSub(user string) (users []entities.User) {
	users, err := b.stor.ListAllAvailableSub(context.Background(), user)
	if err != nil {
		zap.S().Errorln(err)
		return []entities.User{}
	}
	return
}

// List all subsription.
func (b Bot) ListCurrentSub(user string) (users []entities.User) {
	users, err := b.stor.ListCurrentSub(context.Background(), user)
	if err != nil {
		zap.S().Errorln(err)
		return []entities.User{}
	}
	return
}

// Subscribe to user.
func (b Bot) AddSubscription(user, subscribed string, chatID int64) string {
	err := b.stor.AddSubscription(context.Background(), user, subscribed, chatID)
	if err != nil {
		zap.S().Errorln(err)
		return err.Error()
	}
	return fmt.Sprintf("Пользователь %s подписался на  %s успешно. ", user, subscribed)
}

// Unsubscribe from user.
func (b Bot) RemoveSubscription(user, subscribed string) string {
	err := b.stor.RemoveSubscription(context.Background(), user, subscribed)
	if err != nil {
		zap.S().Errorln(err)
		return err.Error()
	}
	return fmt.Sprintf("Пользователь %s отписался от  %s успешно. ", user, subscribed)
}

// Process Callback data.
func (b Bot) Resiver(eAction string) string {
	// Decode action.
	action := DecodeAction(eAction)

	switch action.Type {
	// Make subscribtion.
	case entities.SUB:
		return b.AddSubscription(action.Subscriber, action.Subscribed, action.ChatID)
	case entities.UNSUB:
		return b.RemoveSubscription(action.Subscriber, action.Subscribed)
	}
	return fmt.Sprintln("Subscribe to the user", action.Subscribed)
}

// Function for working with actions.
func EndcodeAction(action entities.Action) string {
	ecnoded, err := json.Marshal(action)
	if err != nil {
		zap.S().Errorln(err)
	}
	return string(ecnoded)
}

func DecodeAction(encoded string) *entities.Action {
	action := &entities.Action{}
	err := json.Unmarshal([]byte(encoded), action)
	if err != nil {
		zap.S().Errorln(err)
	}
	return action
}
