package service

import (
	"sync"

	"github.com/fox-one/fox-notifier/session"
)

// Service service interface
type Service interface {
	Do(s *session.Session)
}

// Run run
func Run(s *session.Session) error {
	group := sync.WaitGroup{}

	notifier, err := createNotifierService(s, s.Notifier())
	if err != nil {
		return err
	}

	services := []Service{
		notifier,
	}

	for idx := range services {
		group.Add(1)
		srv := services[idx]
		go func() {
			srv.Do(s)
			group.Done()
		}()
	}

	group.Wait()
	return nil
}
