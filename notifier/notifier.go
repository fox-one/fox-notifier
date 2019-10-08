package notifier

import (
	"fmt"

	"github.com/fox-one/gin-contrib/session"
	"github.com/fox-one/mixin-sdk/messenger"
	"github.com/fox-one/mixin-sdk/utils"
)

// Notifier notifier
type Notifier struct {
	*messenger.Messenger
}

// ConversationIDFromTopic conversation id from topic
func (n *Notifier) ConversationIDFromTopic(topic string) string {
	return utils.UUIDWithString(fmt.Sprintf("conversation:%s;%s", n.UserID, topic))
}

// FetchOrCreateConversation fetch or create conversation
func (n *Notifier) FetchOrCreateConversation(s *session.Session, topic string) (string, error) {
	conversationID := n.ConversationIDFromTopic(topic)
	if _, f := conversations[conversationID]; f {
		return conversationID, nil
	}

	participants := []*messenger.Participant{
		&messenger.Participant{
			UserID: adminID,
			Role:   "ADMIN",
		},
	}
	if _, err := n.CreateConversation(s.Context(), "GROUP", conversationID, topic, "", "", "", participants); err != nil {
		return "", err
	}

	conversations[conversationID] = true
	t := Topic{
		Topic:          topic,
		ConversationID: conversationID,
	}
	err := s.MysqlWrite().Where("conversation_id = ?", conversationID).FirstOrCreate(&t).Error
	return conversationID, err
}
