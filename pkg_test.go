package log_test

import (
	"errors"
	"testing"

	"github.com/sayden/log"
	"github.com/sayden/log/handlers/memory"
	"github.com/stretchr/testify/assert"
)

type Pet struct {
	Name string
	Age  int
}

func (p *Pet) Fields() log.Fields {
	return log.Fields{
		"name": p.Name,
		"age":  p.Age,
	}
}

func TestInfo(t *testing.T) {
	h := memory.New()
	log.SetHandler(h)

	log.Infof("logged in %s", "Tobi")

	e := h.Entries[0]
	assert.Equal(t, e.GetMessage(), "logged in Tobi")
	assert.Equal(t, e.GetLevel(), log.LevelInfo)
}

func TestFielder(t *testing.T) {
	h := memory.New()
	log.SetHandler(h)

	pet := &Pet{"Tobi", 3}
	log.WithFields(pet).Info("add pet")

	e := h.Entries[0]
	assert.Equal(t, log.Fields{"name": "Tobi", "age": 3}, e.GetFields())
}

// Unstructured logging is supported, but not recommended since it is hard to query.
func Example_unstructured() {
	log.Infof("%s logged in", "Tobi")
}

// Structured logging is supported with fields, and is recommended over the formatted message variants.
func Example_structured() {
	log.WithField("user", "Tobo").Info("logged in")
}

// Errors are passed to WithError(), populating the "error" field.
func Example_errors() {
	err := errors.New("boom")
	log.WithError(err).Error("upload failed")
}

// Multiple fields can be set, via chaining, or WithFields().
func Example_multipleFields() {
	log.WithFields(log.Fields{
		"user": "Tobi",
		"file": "sloth.png",
		"type": "image/png",
	}).Info("upload")
}

// Trace can be used to simplify logging of start and completion events,
// for example an upload which may fail.
func Example_trace() (err error) {
	defer log.Trace("upload").Stop(&err)
	return nil
}
