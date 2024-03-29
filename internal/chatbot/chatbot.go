package chatbot

import (
	"context"
	"time"

	"github.com/evergarden0412/gptea-api/internal"
	"github.com/sashabaranov/go-openai"
)

type Chatbot struct {
	client *openai.Client
}

func New(client *openai.Client) *Chatbot {
	return &Chatbot{
		client: client,
	}
}

func (c *Chatbot) SendChat(ctx context.Context, chatID string, history []*internal.MessageWithScrap, newmsg string) (in, out *internal.Message, err error) {
	lastSeq := 0
	if len(history) != 0 {
		lastSeq = history[0].Seq
	}
	in = &internal.Message{
		ChatID:    chatID,
		Seq:       lastSeq + 1,
		Content:   newmsg,
		CreatedAt: time.Now().UTC(),
		Role:      openai.ChatMessageRoleUser,
	}

	messages := buildMessages(history, in)
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo0613,
		MaxTokens: 1000,
		Messages:  messages,
	}
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	out = &internal.Message{
		ChatID:    chatID,
		Seq:       lastSeq + 2,
		Content:   resp.Choices[0].Message.Content,
		CreatedAt: time.Now().UTC(),
		Role:      resp.Choices[0].Message.Role,
	}
	return
}

// assumes history is sorted in ascending time
func buildMessages(history []*internal.MessageWithScrap, new *internal.Message) []openai.ChatCompletionMessage {
	var res []openai.ChatCompletionMessage
	for _, hist := range history {
		res = append(res, openai.ChatCompletionMessage{Role: hist.Role, Content: hist.Content})
	}
	res = append(res, openai.ChatCompletionMessage{Role: new.Role, Content: new.Content})
	return res
}

func GetSystemMessageRole() string {
	return openai.ChatMessageRoleSystem
}
