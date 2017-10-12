package main

import (
	"os"
	"time"

	"github.com/sayden/log"
	"github.com/sayden/log/handlers/kinesis"
	"github.com/sayden/log/handlers/multi"
	"github.com/sayden/log/handlers/text"
)

func main() {
	log.SetHandler(multi.New(
		text.New(os.Stderr),
		kinesis.New("logs"),
	))

	ctx := log.WithFields(log.Fields{
		"file": "something.png",
		"type": "image/png",
		"user": "tobi",
	})

	for range time.Tick(time.Millisecond * 100) {
		ctx.Info("upload")
		ctx.Info("upload complete")
	}
}
