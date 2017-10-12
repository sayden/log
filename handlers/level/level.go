// Package level implements a level filter handler.
package level

import "github.com/sayden/log"

// Handler implementation.
type Handler struct {
	Level   log.Level
	Handler log.Handler
}

// New handler.
func New(h log.Handler, level log.Level) *Handler {
	return &Handler{
		Level:   level,
		Handler: h,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e log.Interface) error {
	if e.GetLevel() < h.Level {
		return nil
	}

	return h.Handler.HandleLog(e)
}
