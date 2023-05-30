package internal

import (
	"crypto/rand"
	"encoding/base32"
	"time"
)

type Chat struct {
	ID        string     `json:"id" example:"Hjejwerhj"`
	Name      string     `json:"name" example:"basic"`
	CreatedAt *time.Time `json:"createdAt" example:"2021-01-01T00:00:00Z"`
}

func NewChat() (*Chat, error) {
	id := make([]byte, 15) // base32 encoding muiltiple of 5
	_, err := rand.Read(id)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &Chat{
		ID:        base32.StdEncoding.EncodeToString(id),
		CreatedAt: &now,
	}, nil
}

// pk is (chat_id, seq)
type Message struct {
	ChatID    string     `json:"chatID" example:"Hjejwerhj"`
	Seq       int        `json:"seq" example:"1"` // seq starts from 1
	Content   string     `json:"content"`
	Role      string     `json:"role"`
	CreatedAt *time.Time `json:"createdAt"`
}

type Scrapbook struct {
	ID        string    `json:"id" example:"Hjejwerhj"`
	Name      string    `json:"name" example:"basic"`
	CreatedAt time.Time `json:"createdAt" example:"2021-01-01T00:00:00Z"`
}

func (s *Scrapbook) Assign() error {
	id, err := NewUserID()
	if err != nil {
		return err
	}
	s.ID = id
	s.CreatedAt = time.Now().UTC()
	return nil
}

type Scrap struct {
	Memo      string     `json:"memo" example:"hello"`
	CreatedAt *time.Time `json:"createdAt" example:"2021-01-01T00:00:00Z"`
}

func NewUserID() (string, error) {
	id := make([]byte, 15) // base32 encoding muiltiple of 5
	_, err := rand.Read(id)
	if err != nil {
		return "", err
	}
	userID := base32.StdEncoding.EncodeToString(id)
	return userID, nil
}
