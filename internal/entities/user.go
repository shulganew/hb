package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	Tg       string    `db:"tg_user"` // Telegram user nik name ie @user.
	Name     string    `db:"name"`
	PassHash string    `db:"password_hash"`
	Hb       time.Time `db:"hb"` // Happy birthday day.
	UUID     uuid.UUID `db:"user_id"`
}
