package server

import (
	"github.com/fox-one/fox-notifier/session"
	"github.com/gin-gonic/gin"
)

const (
	foxSessionContextKey = "fox.session.context.key"
)

func bindSession(c *gin.Context, s *session.Session) {
	c.Set(foxSessionContextKey, s)
}

// Session session
func Session(c *gin.Context) *session.Session {
	return c.MustGet(foxSessionContextKey).(*session.Session)
}
