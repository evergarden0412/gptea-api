package internal

import "time"

type Chat struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"createdAt"`
}

// pk: ChatID + Seq
type Message struct {
	ChatID    string     `json:"chatID"`
	Seq       int        `json:"seq"`
	Content   string     `json:"content"`
	Role      string     `json:"-"`
	CreatedAt *time.Time `json:"createdAt"`
}

type Scrapbook struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userID"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"createdAt"`
}

type Scrap struct {
	ScrapbookID string     `json:"scrapbookID"`
	MessageID   string     `json:"messageID"`
	CreatedAt   *time.Time `json:"createdAt"`
}
