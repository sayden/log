package level_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sayden/log"
	"github.com/sayden/log/handlers/level"
	"github.com/sayden/log/handlers/memory"
)

func Test(t *testing.T) {
	h := memory.New()

	ctx := log.Logger{
		Handler: level.New(h, log.LevelError),
		Level:   log.LevelInfo,
	}

	ctx.Info("hello")
	ctx.Info("world")
	ctx.Error("boom")

	assert.Len(t, h.Entries, 1)
	assert.Equal(t, h.Entries[0].GetMessage(), "boom")
}
