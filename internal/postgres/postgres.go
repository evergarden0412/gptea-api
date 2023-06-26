package postgres

import (
	"context"
	"database/sql"
	"errors"
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
	ErrUnauthorized     = errors.New("unauthorized")
	ErrChatNotPatched   = errors.New("chat not patched")
	ErrInsertScrapbook  = errors.New("insert scrapbook failed")
	ErrBadScrapbookName = errors.New("bad scrapbook name")
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
	query = `INSERT INTO scrapbooks (id, user_id, name, is_default, created_at) VALUES ($1, $2, $3, $4, $5)`
	scrapbookID, err := internal.NewID()
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, query, scrapbookID, inp.UserID, internal.DefaultScrapbookName, true, inp.CreatedAt); err != nil {
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

func (db *DB) Logout(ctx context.Context, userID string) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	res, err := db.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrUnauthorized
	}
	return nil
}

func (db *DB) Resign(ctx context.Context, userID string) error {
	query := `DELETE FROM users WHERE id = $1`
	res, err := db.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrUnauthorized
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
	res, err := db.db.ExecContext(ctx, query, chat.Name, chat.ID)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrChatNotPatched
	}
	return nil
}

func (db *DB) DeleteChat(ctx context.Context, userID, chatID string) error {
	query := `DELETE FROM chats WHERE id = $1 AND user_id = $2`
	res, err := db.db.ExecContext(ctx, query, chatID, userID)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrUnauthorized
	}
	return nil
}

func (db *DB) GetMyMessages(ctx context.Context, userID, chatID string) ([]*internal.Message, error) {
	_, err := db.SelectMyChat(ctx, userID, chatID)
	if err != nil {
		return nil, err
	}
	query := `SELECT m.chat_id, m.seq, m.content, m.role, m.created_at, COALESCE(s.id, ''), COALESCE(s.memo, ''), COALESCE(s.created_at, '1970-01-01T00:00:00Z') 
		FROM messages AS m
		LEFT JOIN scraps AS s
		ON s.message_chat_id = m.chat_id AND s.message_seq = m.seq
		WHERE m.chat_id = $1
		ORDER BY m.seq DESC`
	rows, err := db.db.QueryContext(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*internal.Message
	for rows.Next() {
		var msg internal.Message
		var scrap internal.Scrap
		if err := rows.Scan(&msg.ChatID, &msg.Seq, &msg.Content, &msg.Role, &msg.CreatedAt, &scrap.ID, &scrap.Memo, &scrap.CreatedAt); err != nil {
			return nil, err
		}
		if scrap.ID != "" {
			msg.Scrap = &scrap
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

func (db *DB) SelectMyScrapbooks(ctx context.Context, userID string) ([]internal.Scrapbook, error) {
	query := `SELECT id, name, is_default, created_at FROM scrapbooks WHERE user_id = $1 ORDER BY created_at ASC`
	rows, err := db.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scrapbooks []internal.Scrapbook
	for rows.Next() {
		var scrapbook internal.Scrapbook
		if err := rows.Scan(&scrapbook.ID, &scrapbook.Name, &scrapbook.IsDefault, &scrapbook.CreatedAt); err != nil {
			return nil, err
		}
		scrapbooks = append(scrapbooks, scrapbook)
	}
	return scrapbooks, nil
}

func (db *DB) InsertScrapbook(ctx context.Context, userID string, inp internal.Scrapbook) error {
	query := `INSERT INTO scrapbooks (id, user_id, name, created_at) VALUES ($1, $2, $3, $4)`
	res, err := db.db.ExecContext(ctx, query, inp.ID, userID, inp.Name, inp.CreatedAt)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrInsertScrapbook
	}
	return nil
}

func (db *DB) DeleteScrapbook(ctx context.Context, userID, scrapbookID string) error {
	query := `DELETE FROM scrapbooks WHERE id = $1 AND user_id = $2 AND is_default = false`
	res, err := db.db.ExecContext(ctx, query, scrapbookID, userID)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrUnauthorized
	}
	return nil
}

func (db *DB) PatchScrapbook(ctx context.Context, userID, scrapbookID, name string) error {
	query := `UPDATE scrapbooks SET name = $1 WHERE id = $2 AND user_id = $3 AND is_default = false`
	res, err := db.db.ExecContext(ctx, query, name, scrapbookID, userID)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrUnauthorized
	}
	return nil
}

func (db *DB) SelectScrapsOnScrapbook(ctx context.Context, userID, scrapbookID string) ([]internal.Scrap, error) {
	query := `SELECT s.id, s.memo, s.created_at, m.chat_id, m.seq, m.content, m.role, m.created_at
		FROM scraps AS s
		INNER JOIN messages AS m
		ON s.message_chat_id = m.chat_id AND s.message_seq = m.seq 
		INNER JOIN scraps_scrapbooks AS ss 
		ON s.id = ss.scrap_id
		INNER JOIN scrapbooks AS sb
		ON ss.scrapbook_id = sb.id
		WHERE sb.user_id = $1 AND sb.id = $2
		ORDER BY s.created_at DESC`
	rows, err := db.db.QueryContext(ctx, query, userID, scrapbookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scraps []internal.Scrap
	for rows.Next() {
		var scrap internal.Scrap
		var msg internal.Message
		if err := rows.Scan(&scrap.ID, &scrap.Memo, &scrap.CreatedAt, &msg.ChatID, &msg.Seq, &msg.Content, &msg.Role, &msg.CreatedAt); err != nil {
			return nil, err
		}
		scrap.Message = &msg
		scraps = append(scraps, scrap)
	}
	return scraps, nil
}

func (db *DB) SelectMyScraps(ctx context.Context, userID string) ([]internal.Scrap, error) {
	query := `SELECT s.id, s.memo, s.created_at, m.chat_id, m.seq, m.content, m.role, m.created_at
		FROM scraps AS s
		INNER JOIN messages AS m
		ON s.message_chat_id = m.chat_id AND s.message_seq = m.seq 
		INNER JOIN chats AS c
		ON m.chat_id = c.id
		WHERE c.user_id = $1
		ORDER BY s.created_at DESC`
	rows, err := db.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scraps []internal.Scrap
	for rows.Next() {
		var scrap internal.Scrap
		var msg internal.Message
		if err := rows.Scan(&scrap.ID, &scrap.Memo, &scrap.CreatedAt, &msg.ChatID, &msg.Seq, &msg.Content, &msg.Role, &msg.CreatedAt); err != nil {
			return nil, err
		}
		scrap.Message = &msg
		scraps = append(scraps, scrap)
	}
	return scraps, nil
}

func (db *DB) InsertScrap(ctx context.Context, userID string, inp internal.Scrap, scrapbookIDs []string) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO scraps (id, memo, created_at, message_chat_id, message_seq)
		SELECT $1, $2, $3, $4, $5
		FROM messages AS m
		INNER JOIN chats AS c ON m.chat_id = c.id
		WHERE m.chat_id = $4 AND m.seq = $5 AND c.user_id = $6`
	res, err := tx.ExecContext(ctx, query, inp.ID, inp.Memo, inp.CreatedAt, inp.Message.ChatID, inp.Message.Seq, userID)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrUnauthorized
	}

	query = `INSERT INTO scraps_scrapbooks (scrap_id, scrapbook_id) 
		SELECT $1, $2
		FROM scrapbooks AS sb
		WHERE sb.user_id = $3 AND sb.id = $2`
	for _, scrapbookID := range scrapbookIDs {
		res, err := tx.ExecContext(ctx, query, inp.ID, scrapbookID, userID)
		if err != nil {
			return err
		}
		if n, err := res.RowsAffected(); err != nil {
			return err
		} else if n == 0 {
			return ErrUnauthorized
		}
	}
	return tx.Commit()
}

func (db *DB) DeleteScrap(ctx context.Context, userID, scrapID string) error {
	query := `DELETE FROM scraps as s
		WHERE s.id = $1 AND 
			(SELECT c.user_id FROM messages AS m
			INNER JOIN chats AS c ON m.chat_id = c.id
			WHERE m.chat_id = s.message_chat_id AND m.seq = s.message_seq) = $2`

	res, err := db.db.ExecContext(ctx, query, scrapID, userID)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrUnauthorized
	}
	return nil
}

func (db *DB) SelectMyScrapbooksOnScrap(ctx context.Context, userID, scrapID string) ([]internal.Scrapbook, error) {
	query := `SELECT sb.id, sb.name, sb.is_default, sb.created_at
		FROM scrapbooks AS sb
		INNER JOIN scraps_scrapbooks AS ss
		ON sb.id = ss.scrapbook_id
		INNER JOIN scraps AS s
		ON ss.scrap_id = s.id
		WHERE s.id = $1 AND sb.user_id = $2`
	rows, err := db.db.QueryContext(ctx, query, scrapID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scrapbooks []internal.Scrapbook
	for rows.Next() {
		var scrapbook internal.Scrapbook
		if err := rows.Scan(&scrapbook.ID, &scrapbook.Name, &scrapbook.IsDefault, &scrapbook.CreatedAt); err != nil {
			return nil, err
		}
		scrapbooks = append(scrapbooks, scrapbook)
	}
	return scrapbooks, nil
}

func (db *DB) InsertScrapOnScrapbook(ctx context.Context, userID, scrapID, scrapbookID string) error {
	query := `INSERT INTO scraps_scrapbooks (scrap_id, scrapbook_id)
		SELECT $1, $2
		FROM scrapbooks AS sb
		WHERE sb.user_id = $3 AND sb.id = $2`
	res, err := db.db.ExecContext(ctx, query, scrapID, scrapbookID, userID)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrUnauthorized
	}
	return nil
}

func (db *DB) DeleteScrapOnScrapbook(ctx context.Context, userID, scrapID, scrapbookID string) error {
	query := `DELETE FROM scraps_scrapbooks AS ss
		WHERE ss.scrap_id = $1 AND ss.scrapbook_id = $2
		AND EXISTS (SELECT 1 FROM scrapbooks AS sb WHERE sb.user_id = $3 AND sb.id = $2)`
	res, err := db.db.ExecContext(ctx, query, scrapID, scrapbookID, userID)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrUnauthorized
	}
	return nil
}
