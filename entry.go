package log

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// assert interface compliance.
var _ Interface = (*Entry)(nil)

// Now returns the current time.
var Now = time.Now

// Entry represents a single log entry.
type Entry struct {
	Logger    *Logger   `json:"-"`
	Fields    Fields    `json:"fields"`
	Level     Level     `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	start     time.Time
	fields    []Fields
	t         Telemetry
}

// NewEntry returns a new entry for `log`.
func NewEntry(log *Logger, t Telemetry) Interface {
	return &Entry{
		Logger: log,
		t:      t,
	}
}

func (e *Entry) Inc(n string, v float64) Interface {
	//first field value
	e.t.Inc(n, v)
	return e
}

func (e *Entry) GetTelemetry() Telemetry {
	return e.t
}

func (e *Entry) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e *Entry) GetLevel() Level {
	return e.Level
}

func (e *Entry) GetMessage() string {
	return e.Message
}

func (e *Entry) GetFields() Fields {
	return e.Fields
}

func (e *Entry) SetMessage(msg string) {
	e.Message = msg
}

func (e *Entry) setStart(t time.Time) {
	e.start = t
}

// WithFields returns a new entry with `fields` set.
func (e *Entry) WithFields(fields Fielder) Interface {
	f := []Fields{}
	f = append(f, e.fields...)
	f = append(f, fields.Fields())
	return &Entry{
		Logger: e.Logger,
		fields: f,
		t:      e.t,
	}
}

// WithField returns a new entry with the `key` and `value` set.
func (e *Entry) WithField(key string, value interface{}) Interface {
	return e.WithFields(Fields{key: value})
}

// WithError returns a new entry with the "error" set to `err`.
//
// The given error may implement .Fielder, if it does the method
// will add all its `.Fields()` into the returned entry.
func (e *Entry) WithError(err error) Interface {
	ctx := e.WithField("error", err.Error())

	if s, ok := err.(stackTracer); ok {
		frame := s.StackTrace()[0]

		name := fmt.Sprintf("%n", frame)
		file := fmt.Sprintf("%+s", frame)
		line := fmt.Sprintf("%d", frame)

		parts := strings.Split(file, "\n\t")
		if len(parts) > 1 {
			file = parts[1]
		}

		ctx = ctx.WithField("source", fmt.Sprintf("%s: %s:%s", name, file, line))
	}

	if f, ok := err.(Fielder); ok {
		ctx = ctx.WithFields(f.Fields())
	}

	return ctx
}

// Debug level message.
func (e *Entry) Debug(msg string) {
	e.Logger.log(LevelDebug, e, msg)
}

// Info level message.
func (e *Entry) Info(msg string) {
	e.Logger.log(LevelInfo, e, msg)
}

// Warn level message.
func (e *Entry) Warn(msg string) {
	e.Logger.log(LevelWarn, e, msg)
}

// Error level message.
func (e *Entry) Error(msg string) {
	e.Logger.log(LevelError, e, msg)
}

// Fatal level message, followed by an exit.
func (e *Entry) Fatal(msg string) {
	e.Logger.log(LevelFatal, e, msg)
	os.Exit(1)
}

// Debugf level formatted message.
func (e *Entry) Debugf(msg string, v ...interface{}) {
	e.Debug(fmt.Sprintf(msg, v...))
}

// Infof level formatted message.
func (e *Entry) Infof(msg string, v ...interface{}) {
	e.Info(fmt.Sprintf(msg, v...))
}

// Warnf level formatted message.
func (e *Entry) Warnf(msg string, v ...interface{}) {
	e.Warn(fmt.Sprintf(msg, v...))
}

// Errorf level formatted message.
func (e *Entry) Errorf(msg string, v ...interface{}) {
	e.Error(fmt.Sprintf(msg, v...))
}

// Fatalf level formatted message, followed by an exit.
func (e *Entry) Fatalf(msg string, v ...interface{}) {
	e.Fatal(fmt.Sprintf(msg, v...))
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func (e *Entry) Trace(msg string) Interface {
	e.Info(msg)
	v := e.WithFields(e.Fields)
	v.SetMessage(msg)
	v.setStart(time.Now())
	return v
}

// Stop should be used with Trace, to fire off the completion message. When
// an `err` is passed the "error" field is set, and the log level is error.
func (e *Entry) Stop(err *error) {
	if err == nil || *err == nil {
		e.WithField("duration", time.Since(e.start)).Info(e.Message)
	} else {
		e.WithField("duration", time.Since(e.start)).WithError(*err).Error(e.Message)
	}
}

// mergedFields returns the fields list collapsed into a single map.
func (e *Entry) mergedFields() Fields {
	f := Fields{}

	for _, fields := range e.fields {
		for k, v := range fields {
			f[k] = v
		}
	}

	return f
}

// finalize returns a copy of the Entry with Fields merged.
func (e *Entry) finalize(level Level, msg string) Interface {
	return &Entry{
		Logger:    e.Logger,
		Fields:    e.mergedFields(),
		Level:     level,
		Message:   msg,
		Timestamp: Now(),
	}
}
