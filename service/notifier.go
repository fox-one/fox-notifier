package service

import (
	"time"

	"github.com/fox-one/fox-notifier/notifier"
	"github.com/fox-one/gin-contrib/session"
	"github.com/smallnest/rpcx/log"
)

const notifierOffsetKey = "notifier:notify_cursor"

type notifierSrv struct {
	notifier *notifier.Notifier
	FromID   int64
}

func createNotifierService(s *session.Session, notifier *notifier.Notifier) (Service, error) {
	srv := &notifierSrv{notifier: notifier}
	from, err := ReadPropertyAsInt64(s, notifierOffsetKey)
	if err != nil {
		return nil, err
	}
	srv.FromID = from
	return srv, nil
}

func (srv *notifierSrv) Do(s *session.Session) {
	duration := time.Millisecond

	for {
		select {
		case <-s.Context().Done():
			return
		case <-time.After(duration):
			if err := srv.do(s); err != nil {
				log.Error("notifier srv", err)
				duration = time.Second
			} else {
				duration = time.Minute
			}
		}
	}
}

func (srv *notifierSrv) do(s *session.Session) error {
	const limit = 100
	msgs, err := notifier.QueryMessages(s, srv.FromID, "", limit)
	if err != nil {
		return err
	}
	if len(msgs) == 0 {
		return nil
	}
	if err := srv.notifier.SendMessages(s, msgs...); err != nil {
		return err
	}
	fromID := msgs[len(msgs)-1].ID
	if err := WriteInt64Property(s, notifierOffsetKey, fromID); err != nil {
		return err
	}
	srv.FromID = fromID
	return nil
}
