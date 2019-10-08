package notifier

import (
	"github.com/fox-one/gin-contrib/session"
	"github.com/jinzhu/gorm"
)

var conversations map[string]bool
var adminID string

// Setup setup package
func Setup(s *session.Session, admin string) error {
	adminID = admin
	var err error
	conversations, err = AllConversations(s)
	return err
}

func setupDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&Topic{}, &Message{}).Error; err != nil {
		return err
	}
	return nil
}

func init() {
	session.RegisterSetdb(setupDB)
}
