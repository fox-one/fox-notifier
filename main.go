package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fox-one/fox-notifier/service"

	"github.com/fox-one/gin-contrib/session"

	"github.com/fox-one/fox-etf/fund"
	"github.com/fox-one/fox-notifier/notifier"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	// NAME name
	NAME = "etf"
	// VERSION version
	VERSION = "null"
	// BUILD build
	BUILD = "null"
)

func main() {
	app := &cli.App{
		Name:        NAME,
		Version:     VERSION + "." + BUILD,
		Description: "Fox Notifier",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "debug, d"},
		},
	}

	s, err := initSession()
	if err != nil {
		panic(err)
	}

	ctx := s.Context()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s = s.WithContext(ctx)
	defer s.Close()

	fund.SetupConfig(s.Session)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Debug("quit app...")
		cancel()

		select {
		case <-time.After(time.Second * 3):
			log.Fatal("quit app timeout")
		}
	}()

	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		notifier.Setup(s.Session, s.Admin())
		return nil
	}

	app.Commands = append(app.Commands, cli.Command{
		Name: "hc",
		Action: func(c *cli.Context) error {
			if err := s.Redis().Ping().Err(); err != nil {
				return err
			}

			if err := s.MysqlRead().DB().Ping(); err != nil {
				return err
			}

			if err := s.MysqlWrite().DB().Ping(); err != nil {
				return err
			}

			return nil
		},
	})

	app.Commands = append(app.Commands, cli.Command{
		Name: "setdb",
		Action: func(c *cli.Context) error {
			return session.Setdb(s.Session)
		},
	})

	app.Commands = append(app.Commands, cli.Command{
		Name: "service",
		Action: func(c *cli.Context) error {
			return service.Run(s)
		},
	})

	// app.Commands = append(app.Commands, &cli.Command{
	// 	Name: "server",
	// 	Flags: []cli.Flag{
	// 		&cli.IntFlag{Name: "port, p", Value: 8081},
	// 	},
	// 	Action: func(c *cli.Context) error {
	// 		return server.Run(s, &server.Option{
	// 			Port:    c.Int("port"),
	// 			Debug:   c.Bool("debug"),
	// 			Version: app.Version,
	// 		})
	// 	},
	// })

	if err := app.Run(os.Args); err != nil {
		log.Errorf("app exit with error: %s", err)
		os.Exit(1)
	}
}
