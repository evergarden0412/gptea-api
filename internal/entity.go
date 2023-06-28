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
	ChatID    string    `json:"chatID" example:"Hjejwerhj"`
	Seq       int       `json:"seq" example:"1"` // seq starts from 1
	Content   string    `json:"content"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type MessageWithScrap struct {
	ChatID    string    `json:"chatID" example:"Hjejwerhj"`
	Seq       int       `json:"seq" example:"1"` // seq starts from 1
	Content   string    `json:"content"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`

	Scrap *Scrap `json:"scrap,omitempty"`
}

type Scrapbook struct {
	ID        string    `json:"id" example:"Hjejwerhj"`
	Name      string    `json:"name" example:"basic"`
	IsDefault bool      `json:"isDefault" example:"true"`
	CreatedAt time.Time `json:"createdAt" example:"2021-01-01T00:00:00Z"`
}

const DefaultScrapbookName = "기본 스크랩북"

func (s *Scrapbook) Assign() error {
	id, err := NewID()
	if err != nil {
		return err
	}
	s.ID = id
	s.CreatedAt = time.Now().UTC()
	return nil
}

type Scrap struct {
	ID        string    `json:"id" example:"Hjejwerhj"`
	Memo      string    `json:"memo" example:"hello"`
	CreatedAt time.Time `json:"createdAt" example:"2021-01-01T00:00:00Z"`
}

func (s *Scrap) Assign() error {
	id, err := NewID()
	if err != nil {
		return err
	}
	s.ID = id
	s.CreatedAt = time.Now().UTC()
	return nil
}

type ScrapWithMessage struct {
	ID        string    `json:"id" example:"Hjejwerhj"`
	Memo      string    `json:"memo" example:"hello"`
	CreatedAt time.Time `json:"createdAt" example:"2021-01-01T00:00:00Z"`

	Message *Message `json:"message,omitempty"`
}

func NewID() (string, error) {
	id := make([]byte, 15) // base32 encoding muiltiple of 5
	_, err := rand.Read(id)
	if err != nil {
		return "", err
	}
	createdID := base32.StdEncoding.EncodeToString(id)
	return createdID, nil
}
