package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/shulganew/hb.git/internal/entities"
)

func (r *Repo) AddSubscription(ctx context.Context, tguser, subscribed string, chatID int64) (err error) {
	query := `
	INSERT INTO subscription (tg_user, subscribed, chat_id) 
	VALUES ($1, $2, $3)
	`

	_, err = r.db.ExecContext(ctx, query, tguser, subscribed, chatID)
	if err != nil {
		var pgErr *pq.Error
		// If exist in DataBase.
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return errors.New("already subscribed")
		}
		return fmt.Errorf("error adding user to Storage: %w", err)
	}
	return
}

func (r *Repo) RemoveSubscription(ctx context.Context, tguser, subscribed string) (err error) {
	query := `
	DELETE FROM subscription WHERE tg_user=$1 AND subscribed=$2;
	`

	_, err = r.db.ExecContext(ctx, query, tguser, subscribed)
	if err != nil {
		return fmt.Errorf("error adding user to Storage: %w", err)
	}
	return
}

// Return users available for sub.
func (r *Repo) ListAllAvailableSub(ctx context.Context, subscriber string) (users []entities.User, err error) {
	query := `SELECT users.user_id, users.tg_user, users.name, users.password_hash, users.hb FROM users LEFT JOIN subscription ON users.tg_user=subscription.tg_user WHERE NOT users.tg_user=$1 AND users.tg_user NOT IN (SELECT subscribed FROM subscription WHERE tg_user=$1)`
	users = []entities.User{}
	err = r.db.SelectContext(ctx, &users, query, subscriber)
	if err != nil {
		return nil, fmt.Errorf("error during get user by login from storage. User not valid: %w", err)
	}
	return users, nil
}

// Return sub for user.
func (r *Repo) ListCurrentSub(ctx context.Context, subscriber string) (users []entities.User, err error) {
	query := `SELECT * FROM users WHERE tg_user IN (SELECT subscribed FROM subscription Where tg_user=$1)`
	users = []entities.User{}
	err = r.db.SelectContext(ctx, &users, query, subscriber)
	if err != nil {
		return nil, fmt.Errorf("error during get user by login from storage. User not valid: %w", err)
	}
	return users, nil
}
