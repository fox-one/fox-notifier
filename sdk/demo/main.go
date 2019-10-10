package main

import (
	"context"
	"log"

	notifier "github.com/fox-one/fox-notifier/sdk"
)

func main() {
	n := notifier.NewNotifier("http://localhost:8888")
	msg, err := n.NotifyMessage(context.TODO(), "", "test", "test")
	log.Println(msg, err)
}
