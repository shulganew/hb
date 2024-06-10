package services

import (
	"context"

	"github.com/shulganew/hb.git/internal/api/oapi"
	"github.com/shulganew/hb.git/internal/config"
	"github.com/shulganew/hb.git/internal/entities"
)

// User creation, registration
type Happy struct {
	stor   Happyer
	conf   config.Config
	memory Memer
}

func NewHappy(ctx context.Context, stor Happyer, conf config.Config, m Memer) *Happy {
	hb := &Happy{stor: stor, conf: conf, memory: m}
	return hb
}

type Happyer interface {
	// Happy user.
	AddUser(ctx context.Context, tguser, name, pwhash, hbd string) (err error)

	// Happy bot.
	ListAll(ctx context.Context) ([]entities.User, error)
	ListAllAvailableSub(ctx context.Context, user string) ([]entities.User, error)
	ListCurrentSub(ctx context.Context, user string) ([]entities.User, error)
	AddSubscription(ctx context.Context, tguser, subscribed string, chatID int64) (err error)
	RemoveSubscription(ctx context.Context, tguser, subscribed string) (err error)
	GetNotifyChats(ctx context.Context, tguser string) (chatIDs []int64, err error)
	// Cron.
	// Return happy users.
	ListBithdayMan(ctx context.Context) ([]entities.User, error)

	// Entities credentials methods (site, card, text, file)
}
type Memer interface {
	// Memory storage.
	Add(user string)
	Get(user string) bool
}

// Check interfaces.
var _ oapi.ServerInterface = (*Happy)(nil)
