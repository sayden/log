package log

import (
	"bytes"
	"errors"
	"strings"
)

// ErrInvalidLevel is returned if the severity level is invalid.
var ErrInvalidLevel = errors.New("invalid level")

// Level of severity.
type Level int

// Log levels.
const (
	InvalidLevel Level = iota - 1
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levelNames = [...]string{
	LevelDebug: "debug",
	LevelInfo:  "info",
	LevelWarn:  "warn",
	LevelError: "error",
	LevelFatal: "fatal",
}

var levelStrings = map[string]Level{
	"debug":   LevelDebug,
	"info":    LevelInfo,
	"warn":    LevelWarn,
	"warning": LevelWarn,
	"error":   LevelError,
	"fatal":   LevelFatal,
}

// String implementation.
func (l Level) String() string {
	return levelNames[l]
}

// MarshalJSON implementation.
func (l Level) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

// UnmarshalJSON implementation.
func (l *Level) UnmarshalJSON(b []byte) error {
	v, err := ParseLevel(string(bytes.Trim(b, `"`)))
	if err != nil {
		return err
	}

	*l = v
	return nil
}

// ParseLevel parses level string.
func ParseLevel(s string) (Level, error) {
	l, ok := levelStrings[strings.ToLower(s)]
	if !ok {
		return InvalidLevel, ErrInvalidLevel
	}

	return l, nil
}

// MustParseLevel parses level string or panics.
func MustParseLevel(s string) Level {
	l, err := ParseLevel(s)
	if err != nil {
		panic("invalid log level")
	}

	return l
}
