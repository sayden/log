package log

import (
	stdlog "log"
	"sort"
	"time"
)

// assert interface compliance.
var _ Interface = (*Logger)(nil)

// Fielder is an interface for providing fields to custom types.
type Fielder interface {
	Fields() Fields
}

// Fields represents a map of entry level data used for structured logging.
type Fields map[string]interface{}

// Fields implements Fielder.
func (f Fields) Fields() Fields {
	return f
}

// Get field value by name.
func (f Fields) Get(name string) interface{} {
	return f[name]
}

// Names returns field names sorted.
func (f Fields) Names() (v []string) {
	for k := range f {
		v = append(v, k)
	}

	sort.Strings(v)
	return
}

// The HandlerFunc type is an adapter to allow the use of ordinary functions as
// log handlers. If f is a function with the appropriate signature,
// HandlerFunc(f) is a Handler object that calls f.
type HandlerFunc func(Interface) error

// HandleLog calls f(e).
func (f HandlerFunc) HandleLog(e Interface) error {
	return f(e)
}

// Handler is used to handle log events, outputting them to
// stdio or sending them to remote services. See the "handlers"
// directory for implementations.
//
// It is left up to Handlers to implement thread-safety.
type Handler interface {
	HandleLog(Interface) error
}

// Logger represents a logger with configurable Level and Handler.
type Logger struct {
	Handler Handler
	Level   Level
	t       Telemetry
}

func (e *Logger) Inc(n string, v float64) Interface {
	//first field value
	e.t.Inc(n,v)
	return e
}

func (e *Logger) GetTelemetry() Telemetry {
	return e.t
}

func (e *Logger) Stop(_ *error) {
	Error("Stop Does nothing")
}

func (e *Logger) GetLevel() Level {
	return e.Level
}

func (e *Logger) GetMessage() string {
	Error("GetMessage Does nothing")
	return ""
}

func (e *Logger) GetFields() Fields {
	Error("GetFields Does nothing")
	return nil
}

func (e *Logger) GetTimestamp() time.Time {
	Error("GetTimestamp Does nothing")
	return time.Time{}
}

func (e *Logger) finalize(_ Level, _ string) Interface {
	Error("finalize Does nothing")
	return nil
}

func (e *Logger) mergedFields() Fields {
	Error("mergedFields Does nothing")
	return nil
}

func (e *Logger) SetMessage(msg string) {
	Error("SetMessageDoes nothing")
}
func (l *Logger) setStart(t time.Time) {
	Error("setStart Does nothing")
}

// WithFields returns a new entry with `fields` set.
func (l *Logger) WithFields(fields Fielder) Interface {
	return NewEntry(l, l.t).WithFields(fields.Fields())
}

// WithField returns a new entry with the `key` and `value` set.
//
// Note that the `key` should not have spaces in it - use camel
// case or underscores
func (l *Logger) WithField(key string, value interface{}) Interface {
	return NewEntry(l, l.t).WithField(key, value)
}

// WithError returns a new entry with the "error" set to `err`.
func (l *Logger) WithError(err error) Interface {
	return NewEntry(l, l.t).WithError(err)
}

// Debug level message.
func (l *Logger) Debug(msg string) {
	NewEntry(l, l.t).Debug(msg)
}

// Info level message.
func (l *Logger) Info(msg string) {
	NewEntry(l, l.t).Info(msg)
}

// Warn level message.
func (l *Logger) Warn(msg string) {
	NewEntry(l, l.t).Warn(msg)
}

// Error level message.
func (l *Logger) Error(msg string) {
	NewEntry(l, l.t).Error(msg)
}

// Fatal level message, followed by an exit.
func (l *Logger) Fatal(msg string) {
	NewEntry(l, l.t).Fatal(msg)
}

// Debugf level formatted message.
func (l *Logger) Debugf(msg string, v ...interface{}) {
	NewEntry(l, l.t).Debugf(msg, v...)
}

// Infof level formatted message.
func (l *Logger) Infof(msg string, v ...interface{}) {
	NewEntry(l, l.t).Infof(msg, v...)
}

// Warnf level formatted message.
func (l *Logger) Warnf(msg string, v ...interface{}) {
	NewEntry(l, l.t).Warnf(msg, v...)
}

// Errorf level formatted message.
func (l *Logger) Errorf(msg string, v ...interface{}) {
	NewEntry(l, l.t).Errorf(msg, v...)
}

// Fatalf level formatted message, followed by an exit.
func (l *Logger) Fatalf(msg string, v ...interface{}) {
	NewEntry(l, l.t).Fatalf(msg, v...)
}

// Trace returns a new entry with a Stop method to fire off
// a corresponding completion log, useful with defer.
func (l *Logger) Trace(msg string) Interface {
	return NewEntry(l, l.t).Trace(msg)
}

// log the message, invoking the handler. We clone the entry here
// to bypass the overhead in Entry methods when the level is not
// met.
func (l *Logger) log(level Level, e *Entry, msg string) {
	if level < l.Level {
		return
	}

	if err := l.Handler.HandleLog(e.finalize(level, msg)); err != nil {
		stdlog.Printf("error logging: %s", err)
	}
}
