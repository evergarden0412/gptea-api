package internal

import "time"

type Chat struct {
	ID        string     `json:"id" example:"Hjejwerhj"`
	Name      string     `json:"name" example:"basic"`
	CreatedAt *time.Time `json:"createdAt" example:"2021-01-01T00:00:00Z"`
}

// pk is (chat_id, seq)
type Message struct {
	ChatID    string     `json:"chatID" example:"Hjejwerhj"`
	Seq       int        `json:"seq" example:"1"` // seq starts from 1
	Content   string     `json:"content"`
	Role      string     `json:"-"`
	CreatedAt *time.Time `json:"createdAt"`
}

type Scrapbook struct {
	ID        string     `json:"id" example:"Hjejwerhj"`
	Name      string     `json:"name" example:"basic"`
	CreatedAt *time.Time `json:"createdAt" example:"2021-01-01T00:00:00Z"`
}

type Scrap struct {
	Memo      string     `json:"memo" example:"hello"`
	CreatedAt *time.Time `json:"createdAt" example:"2021-01-01T00:00:00Z"`
}
