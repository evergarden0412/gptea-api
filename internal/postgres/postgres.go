package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/evergarden0412/gptea-api/internal"
)

type DB struct {
	db *sql.DB
}

func New(db *sql.DB) *DB {
	return &DB{db: db}
}

var (
	ErrUnauthorized = fmt.Errorf("unauthorized")
)

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

func (db *DB) Resign(ctx context.Context, userID string) error {
	query := `DELETE FROM users WHERE id = $1`
	if _, err := db.db.ExecContext(ctx, query, userID); err != nil {
		return err
	}
	return nil
}

func (db *DB) IsRefreshTokenExists(ctx context.Context, userID, tokenID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM refresh_tokens WHERE user_id = $1 AND token_id = $2)`
	var exists bool
	if err := db.db.QueryRowContext(ctx, query, userID, tokenID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (db *DB) UpsertRefreshToken(ctx context.Context, userID, tokenID string) error {
	query := `INSERT INTO refresh_tokens (user_id, token_id) VALUES ($1, $2)
			ON CONFLICT (user_id) DO UPDATE SET token_id = $2`
	_, err := db.db.ExecContext(ctx, query, userID, tokenID)
	return err
}

func (db *DB) SelectMyChats(ctx context.Context, userID string) ([]internal.Chat, error) {
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

func (db *DB) SelectMyChat(ctx context.Context, userID, chatID string) (internal.Chat, error) {
	query := `SELECT id, user_id, name, created_at FROM chats WHERE id = $1`
	var chat internal.Chat
	var chatUserID string
	if err := db.db.QueryRowContext(ctx, query, chatID).Scan(&chat.ID, &chatUserID, &chat.Name, &chat.CreatedAt); err != nil {
		return internal.Chat{}, err
	}
	if chatUserID != userID {
		return internal.Chat{}, ErrUnauthorized
	}
	return chat, nil
}

func (db *DB) InsertChat(ctx context.Context, userID string, inp internal.Chat) error {
	query := `INSERT INTO chats (id, user_id, name, created_at) VALUES ($1, $2, $3, $4)`
	if _, err := db.db.ExecContext(ctx, query, inp.ID, userID, inp.Name, inp.CreatedAt); err != nil {
		return err
	}
	return nil
}

func (db *DB) PatchChat(ctx context.Context, userID string, inp internal.Chat) error {
	chat, err := db.SelectMyChat(ctx, userID, inp.ID)
	if err != nil {
		return err
	}
	if inp.Name != "" {
		chat.Name = inp.Name
	}
	query := `UPDATE chats SET name = $1 WHERE id = $2`
	if _, err := db.db.ExecContext(ctx, query, chat.Name, chat.ID); err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteChat(ctx context.Context, userID, chatID string) error {
	query := `DELETE FROM chats WHERE id = $1 AND user_id = $2`
	if _, err := db.db.ExecContext(ctx, query, chatID, userID); err != nil {
		return err
	}
	return nil
}

func (db *DB) GetMyMessages(ctx context.Context, userID, chatID string) ([]*internal.Message, error) {
	_, err := db.SelectMyChat(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}
	query := `SELECT chat_id, seq, content, role, created_at FROM messages WHERE chat_id = $1 ORDER BY seq DESC`
	rows, err := db.db.QueryContext(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*internal.Message
	for rows.Next() {
		var msg internal.Message
		if err := rows.Scan(&msg.ChatID, &msg.Seq, &msg.Content, &msg.Role, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}

func (db *DB) InsertMessage(ctx context.Context, userID string, inp internal.Message) error {
	_, err := db.SelectMyChat(ctx, userID, inp.ChatID)
	if err != nil {
		return err
	}
	query := `INSERT INTO messages (chat_id, seq, content, role, created_at) VALUES ($1, $2, $3, $4, $5)`
	if _, err := db.db.ExecContext(ctx, query, inp.ChatID, inp.Seq, inp.Content, inp.Role, inp.CreatedAt); err != nil {
		return err
	}
	return nil
}
