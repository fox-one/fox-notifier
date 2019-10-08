package notifier

import (
	"time"

	"github.com/fox-one/gin-contrib/session"
)

// Topic topic
type Topic struct {
	ID             int64     `gorm:"PRIMARY_KEY" json:"id"`
	Topic          string    `json:"topic"`
	ConversationID string    `gorm:"SIZE:36;UNIQUE_INDEX;" json:"conversation_id"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
}

// TableName table
func (Topic) TableName() string {
	return "topics"
}

// AllConversations all conversations
func AllConversations(s *session.Session) (map[string]bool, error) {
	var topics []*Topic
	if err := s.MysqlRead().Find(&topics).Error; err != nil {
		return nil, err
	}

	conversations := make(map[string]bool, len(topics))
	for _, t := range topics {
		conversations[t.ConversationID] = true
	}
	return conversations, nil
}
