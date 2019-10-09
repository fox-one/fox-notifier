package notifier

import (
	"context"
	"errors"

	"github.com/fox-one/httpclient"
	jsoniter "github.com/json-iterator/go"
)

// Notifier notifier
type Notifier struct {
	client *httpclient.Client
}

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

// NewNotifier new notifier
func NewNotifier(apiBase string) *Notifier {
	return &Notifier{
		client: httpclient.NewClient(apiBase),
	}
}

// NotifyMessage notify message
func (n *Notifier) NotifyMessage(ctx context.Context, messageID, topic, message string) error {
	bts, err := n.client.POST("/message").
		P("message_id", messageID).
		P("topic", topic).
		P("message", message).
		Do(ctx).Bytes()

	if err != nil {
		return err
	}

	var resp struct {
		Code   int    `json:"code"`
		ErrMsg string `json:"msg"`
	}
	if err := jsoniter.Unmarshal(bts, &resp); err != nil {
		return err
	}

	switch resp.Code {
	case 0:
		return nil

	case ErrCodeInvalidInput:
		return ErrInvalidInput
	case ErrCodeServerFault:
		return ErrServerFault

	default:
		return ErrUnknown
	}
}
