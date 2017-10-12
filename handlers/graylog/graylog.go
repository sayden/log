// Package implements a Graylog-backed handler.
package graylog

import (
	"github.com/sayden/log"
	"github.com/aphistic/golf"
)

// Handler implementation.
type Handler struct {
	logger *golf.Logger
	client *golf.Client
}

// New handler.
// Connection string should be in format "udp://<ip_address>:<port>".
// Server should have GELF input enabled on that port.
func New(url string) (*Handler, error) {
	c, err := golf.NewClient()
	if err != nil {
		return nil, err
	}

	err = c.Dial(url)
	if err != nil {
		return nil, err
	}

	l, err := c.NewLogger()
	if err != nil {
		return nil, err
	}

	return &Handler{
		logger: l,
		client: c,
	}, nil
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e log.Interface) error {
	switch e.GetLevel() {
	case log.LevelDebug:
		return h.logger.Dbgm(e.GetFields(), e.GetMessage())
	case log.LevelInfo:
		return h.logger.Infom(e.GetFields(), e.GetMessage())
	case log.LevelWarn:
		return h.logger.Warnm(e.GetFields(), e.GetMessage())
	case log.LevelError:
		return h.logger.Errm(e.GetFields(), e.GetMessage())
	case log.LevelFatal:
		return h.logger.Critm(e.GetFields(), e.GetMessage())
	}

	return nil
}

// Closes connection to server, flushing message queue.
func (h *Handler) Close() error {
	return h.client.Close()
}
