package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/shulganew/hb.git/internal/entities"
)

// Return sub for user.
func (r *Repo) ListBithdayMan(ctx context.Context) (users []entities.User, err error) {
	query := `SELECT * FROM users WHERE EXTRACT(MONTH FROM hb) = $1 and EXTRACT(DAY FROM hb) = $2;`
	users = []entities.User{}
	now := time.Now()
	month := now.Month()
	day := now.Day()

	err = r.db.SelectContext(ctx, &users, query, int(month), day)
	if err != nil {
		return nil, fmt.Errorf("error during get hn users from storage: %w", err)
	}
	return users, nil
}

// Return sub for user.
func (r *Repo) GetNotifyChats(ctx context.Context, subscribed string) (chatIDs []int64, err error) {
	query := `SELECT DISTINCT chat_id FROM subscription where subscribed=$1`

	err = r.db.SelectContext(ctx, &chatIDs, query, subscribed)
	if err != nil {
		return nil, fmt.Errorf("error during get hn users from storage: %w", err)
	}
	return chatIDs, nil
}
