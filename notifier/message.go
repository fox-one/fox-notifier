package notifier

import (
	"time"

	"github.com/fox-one/gin-contrib/session"
)

// Message message
type Message struct {
	ID             int64     `gorm:"PRIMARY_KEY" json:"id"`
	Topic          string    `json:"topic"`
	ConversationID string    `gorm:"SIZE:36;" json:"conversation_id"`
	MessageID      string    `gorm:"SIZE:36;UNIQUE;" json:"message_id"`
	Message        string    `gorm:"TYPE:LONGTEXT;" json:"message"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

// QueryMessages query messages
func QueryMessages(s *session.Session, from int64, conversationID string, limit int) ([]*Message, error) {
	query := s.MysqlRead().Limit(limit)
	if from > 0 {
		query = query.Where("id > ?", from)
	}
	if conversationID != "" {
		query = query.Where("conversation_id = ?", conversationID)
	}
	var msgs []*Message
	err := query.Find(&msgs).Error
	return msgs, err
}
