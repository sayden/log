// Package text implements a development-friendly textual handler.
package text

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/sayden/log"
)

// Default handler outputting to stderr.
var Default = New(os.Stderr)

// start time.
var start = time.Now()

// colors.
const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	gray   = 37
)

// Colors mapping.
var Colors = [...]int{
	log.LevelDebug: gray,
	log.LevelInfo:  blue,
	log.LevelWarn:  yellow,
	log.LevelError: red,
	log.LevelFatal: red,
}

// Strings mapping.
var Strings = [...]string{
	log.LevelDebug: "DEBUG",
	log.LevelInfo:  "INFO",
	log.LevelWarn:  "WARN",
	log.LevelError: "ERROR",
	log.LevelFatal: "FATAL",
}

// Handler implementation.
type Handler struct {
	mu     sync.Mutex
	Writer io.Writer
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		Writer: w,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e log.Interface) error {
	color := Colors[e.GetLevel()]
	level := Strings[e.GetLevel()]
	names := e.GetFields().Names()

	h.mu.Lock()
	defer h.mu.Unlock()

	ts := time.Since(start) / time.Second
	fmt.Fprintf(h.Writer, "\033[%dm%6s\033[0m[%04d] %-25s", color, level, ts, e.GetMessage())

	for _, name := range names {
		fmt.Fprintf(h.Writer, " \033[%dm%s\033[0m=%v", color, name, e.GetFields().Get(name))
	}

	fmt.Fprintln(h.Writer)

	return nil
}
