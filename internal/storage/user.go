package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/shulganew/hb.git/internal/entities"
	"go.uber.org/zap"
)

func (r *Repo) AddUser(ctx context.Context, tguser, name, pwhash, hbd string) (err error) {
	query := `
	INSERT INTO users (tg_user, name, password_hash, hb) 
	VALUES ($1, $2, $3, $4)
	`

	_, err = r.db.ExecContext(ctx, query, tguser, name, pwhash, hbd)
	if err != nil {
		var pgErr *pq.Error
		// If URL exist in DataBase.
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return pgErr
		}
		return fmt.Errorf("error adding user to Storage: %w", err)
	}
	return
}

// Retrive User by login.
func (r *Repo) GetByLogin(ctx context.Context, login string) (*entities.User, error) {
	query := `
	SELECT user_id, password_hash, email, otp_key
	FROM users 
	WHERE login = $1
	`
	user := entities.User{Login: login}
	zap.S().Infoln("user login: ", login)
	err := r.db.GetContext(ctx, &user, query, login)
	if err != nil {
		return nil, fmt.Errorf("error during get user by login from storage. User not valid: %w", err)
	}
	return &user, nil
}
