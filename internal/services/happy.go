package services

import (
	"context"

	"github.com/shulganew/hb.git/internal/api/oapi"
	"github.com/shulganew/hb.git/internal/config"
	"github.com/shulganew/hb.git/internal/entities"
)

// User creation, registration
type Happy struct {
	stor Happyer
	conf config.Config
}

func NewHappy(ctx context.Context, stor Happyer, conf config.Config) *Happy {
	hb := &Happy{stor: stor, conf: conf}

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

	// Entities credentials methods (site, card, text, file)
}

// Check interfaces.
var _ oapi.ServerInterface = (*Happy)(nil)
