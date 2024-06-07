package services

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/shulganew/hb.git/internal/api/oapi"
	"github.com/shulganew/hb.git/internal/config"
	"github.com/shulganew/hb.git/internal/entities"
)

// User creation, registration
type Happy struct {
	stor Happyer

	conf config.Config
	// hashes []entities.EKeyMem // Autentificated telegram users.
}

func NewHappy(ctx context.Context, stor Happyer, conf config.Config) *Happy {
	hb := &Happy{stor: stor, conf: conf}

	return hb
}

type Happyer interface {
	// Happy
	AddUser(ctx context.Context, login, hash, email string) (userID *uuid.UUID, err error)
	GetByLogin(ctx context.Context, login string) (userID *entities.User, err error)

	// Entities credentials methods (site, card, text, file)
}

// Check interfaces.
var _ oapi.ServerInterface = (*Happy)(nil)
