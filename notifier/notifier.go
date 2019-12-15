package notifier

import (
	"encoding/base64"
	"fmt"

	"github.com/fox-one/gin-contrib/session"
	mixin "github.com/fox-one/mixin-sdk"
	"github.com/fox-one/mixin-sdk/utils"
)

// Notifier notifier
type Notifier struct {
	*mixin.User
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

	participants := []*mixin.Participant{
		&mixin.Participant{
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

// CreateMessage create messaage
func (n *Notifier) CreateMessage(s *session.Session, messageID, topic, message string) (*Message, error) {
	conversationID, err := n.FetchOrCreateConversation(s, topic)
	if err != nil {
		return nil, err
	}

	msg := Message{
		ConversationID: conversationID,
		MessageID:      messageID,
		Topic:          topic,
		Message:        message,
	}
	err = s.MysqlWrite().Where("message_id = ?", messageID).FirstOrCreate(&msg).Error
	return &msg, err
}

// SendMessages send messages
func (n *Notifier) SendMessages(s *session.Session, msgs ...*Message) error {
	if len(msgs) == 0 {
		return nil
	}

	var arr = make([]mixin.MessageRequest, len(msgs))
	for idx, msg := range msgs {
		arr[idx] = mixin.MessageRequest{
			ConversationID: msg.ConversationID,
			MessageID:      msg.MessageID,
			Category:       "PLAIN_TEXT",
			Data:           base64.StdEncoding.EncodeToString([]byte(msg.Message)),
		}
	}
	return n.User.SendMessages(s.Context(), arr...)
}
