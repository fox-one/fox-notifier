package notifier

import (
	"context"
	"errors"
	"time"

	"github.com/fox-one/httpclient"
	jsoniter "github.com/json-iterator/go"
)

const (
	// ErrCodeInvalidInput err invalid input
	ErrCodeInvalidInput = 1001
	// ErrCodeServerFault err server fault
	ErrCodeServerFault = 1002
)

var (
	// ErrInvalidInput err invalid input
	ErrInvalidInput = errors.New("invalid input")
	// ErrServerFault err server fault
	ErrServerFault = errors.New("internal server error")

	// ErrUnknown unknonw error
	ErrUnknown = errors.New("unknown error")
)

// Notifier notifier
type Notifier struct {
	client *httpclient.Client
}

// NewNotifier new notifier
func NewNotifier(apiBase string) *Notifier {
	return &Notifier{
		client: httpclient.NewClient(apiBase),
	}
}

// Message message
type Message struct {
	ConversationID string    `json:"conversation_id"`
	MessageID      string    `json:"message_id"`
	Topic          string    `json:"topic"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// NotifyMessage notify message
func (n *Notifier) NotifyMessage(ctx context.Context, messageID, topic, message string) (*Message, error) {
	bts, err := n.client.POST("/message").
		P("message_id", messageID).
		P("topic", topic).
		P("message", message).
		Do(ctx).Bytes()

	if err != nil {
		return nil, err
	}

	var resp struct {
		Code   int    `json:"code"`
		ErrMsg string `json:"msg"`

		Message *Message `json:"data"`
	}
	if err := jsoniter.Unmarshal(bts, &resp); err != nil {
		return nil, err
	}

	switch resp.Code {
	case 0:
		return resp.Message, nil

	case ErrCodeInvalidInput:
		return nil, ErrInvalidInput
	case ErrCodeServerFault:
		return nil, ErrServerFault

	default:
		return nil, ErrUnknown
	}
}
