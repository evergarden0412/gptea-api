package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/evergarden0412/gptea-api/internal"
)

type DB struct {
	db *sql.DB
}

func New(db *sql.DB) *DB {
	return &DB{db: db}
}

type RegisterInput struct {
	UserID         string
	CredentialType string
	CredentialID   string
	CreatedAt      *time.Time
}

func (db *DB) Register(ctx context.Context, inp RegisterInput) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `INSERT INTO users (id, created_at) VALUES ($1, $2)`
	if _, err := tx.ExecContext(ctx, query, inp.UserID, inp.CreatedAt); err != nil {
		return err
	}
	query = `INSERT INTO user_credentials (user_id, credential_type, credential_id) VALUES ($1, $2, $3)`
	if _, err := tx.ExecContext(ctx, query, inp.UserID, inp.CredentialType, inp.CredentialID); err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) SignIn(ctx context.Context, credentialType, credentialID string) (string, error) {
	var userID string
	query := `SELECT user_id FROM user_credentials WHERE credential_type = $1 AND credential_id = $2`
	if err := db.db.QueryRowContext(ctx, query, credentialType, credentialID).Scan(&userID); err != nil {
		return "", err
	}
	return userID, nil
}

func (db *DB) SelectChats(ctx context.Context, userID string) ([]internal.Chat, error) {
	query := `SELECT id, name, created_at FROM chats WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := db.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chats []internal.Chat
	for rows.Next() {
		var chat internal.Chat
		if err := rows.Scan(&chat.ID, &chat.Name, &chat.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

func (db *DB) InsertChat(ctx context.Context, userID string, inp internal.Chat) error {
	query := `INSERT INTO chats (id, user_id, name, created_at) VALUES ($1, $2, $3, $4)`
	if _, err := db.db.ExecContext(ctx, query, inp.ID, userID, inp.Name, inp.CreatedAt); err != nil {
		return err
	}
	return nil
}