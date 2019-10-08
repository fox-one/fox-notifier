package session

import (
	"bytes"
	"context"
	"io"

	"github.com/fox-one/fox-notifier/notifier"
	gSession "github.com/fox-one/gin-contrib/session"
	"github.com/fox-one/mixin-sdk/messenger"
	"github.com/spf13/viper"
)

// Session session
type Session struct {
	*gSession.Session

	admin string

	// shared configuration
	v *viper.Viper

	// mixin dapp
	notifier *notifier.Notifier
}

// New new session with data
func New(data []byte) (*Session, error) {
	return NewWithReader(bytes.NewReader(data))
}

// NewWithReader new session with reader
func NewWithReader(r io.Reader) (*Session, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(r); err != nil {
		return nil, err
	}

	return NewWithViper(v)
}

// NewWithViper new session with viper
func NewWithViper(v *viper.Viper) (*Session, error) {
	s := &Session{
		Session: gSession.NewWithViper(v),
		v:       v,
	}
	s.admin = v.GetString("admin")

	notifier, err := createNotifier(v)
	if err != nil {
		return nil, err
	}
	s.notifier = notifier
	return s, nil
}

// Copy copy
func (s *Session) Copy() *Session {
	return &Session{
		Session:  s.Session.Copy(),
		v:        s.v,
		admin:    s.admin,
		notifier: s.notifier,
	}
}

// MysqlReadOnWrite mysql read on write
func (s *Session) MysqlReadOnWrite() *Session {
	s = s.Copy()
	s.Session = s.Session.MysqlReadOnWrite()
	return s
}

// WithContext with context
func (s *Session) WithContext(ctx context.Context) *Session {
	if ctx == nil {
		panic("nil context")
	}

	cp := s.Copy()
	cp.Session = cp.Session.WithContext(ctx)
	return cp
}

// Notifier notifier
func (s *Session) Notifier() *notifier.Notifier {
	return s.notifier
}

// Admin admin id
func (s *Session) Admin() string {
	return s.admin
}

func createNotifier(v *viper.Viper) (*notifier.Notifier, error) {
	m, err := messenger.NewMessengerWithSession(
		v.GetString("mixin.client_id"),
		v.GetString("mixin.session_id"),
		v.GetString("mixin.session_key"),
	)
	if err != nil {
		return nil, err
	}
	return &notifier.Notifier{Messenger: m}, nil
}
