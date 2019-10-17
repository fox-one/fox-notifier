package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fox-one/fox-notifier/session"
	"github.com/fox-one/gin-contrib/errors"
	"github.com/fox-one/gin-contrib/gin_helper"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

// Option options
type Option struct {
	Version string
	Port    int
	Debug   bool
}

var (
	// ErrInvalidInput err invalid input
	ErrInvalidInput = errors.New(1001, "invalid input")
	// ErrServerFault err server fault
	ErrServerFault = errors.New(1002, "internal server error", http.StatusInternalServerError)
)

// Run run server
func Run(s *session.Session, opt *Option) error {
	if opt.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	loc, _ := time.LoadLocation("Asia/Chongqing")
	r.Use(gin_helper.Log(loc))

	corsOp := cors.DefaultConfig()
	corsOp.AllowCredentials = true
	corsOp.AllowHeaders = []string{
		"Authorization",
		"Origin",
		"Content-Type",
	}
	corsOp.AllowOriginFunc = func(origin string) bool { return true }
	r.Use(cors.New(corsOp))
	r.Use(func(c *gin.Context) { bindSession(c, s) })

	addr := fmt.Sprintf(":%d", opt.Port)

	route(s, &r.RouterGroup, opt)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	select {
	case <-s.Context().Done():
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return err
		}
	}

	log.Debugf("srv shutdone")
	return nil
}

func healthCheck(version string, detail bool) gin.HandlerFunc {
	start := time.Now()

	return func(c *gin.Context) {
		views := gin.H{
			"version":  version,
			"duration": time.Since(start).String(),
		}

		if detail {
			s := Session(c)
			views["mysql_read"] = s.MysqlRead().DB().Ping() == nil
			views["mysql_write"] = s.MysqlWrite().DB().Ping() == nil
		}

		gin_helper.Data(c, views)
	}
}

func createMessage(c *gin.Context) {
	var input struct {
		MessageID string `json:"message_id"`
		Topic     string `json:"topic" binding:"required"`
		Message   string `json:"message" binding:"required"`
	}
	if err := gin_helper.BindJson(c, &input); err != nil {
		gin_helper.FailError(c, ErrInvalidInput, err)
		return
	}
	if input.MessageID == "" {
		input.MessageID = uuid.Must(uuid.NewV4()).String()
	}

	session := Session(c)
	msg, err := session.Notifier().CreateMessage(session.Session, input.MessageID, input.Topic, input.Message)
	if err != nil {
		gin_helper.FailError(c, ErrServerFault, err)
		return
	}

	gin_helper.Data(c, msg)
}

func route(s *session.Session, r *gin.RouterGroup, opt *Option) {
	r.GET("/hc", healthCheck(opt.Version, false))

	r.POST("/message", createMessage)
}
