package ssetest

import (
	"testing"

	"github.com/patdowney/http-sse"
)

func TestEnsureImplementsEventWriter(t *testing.T) {
	r := NewStreamRecorder()
	i := interface{}(r)
	_, ok := i.(sse.EventWriter)
	if !ok {
		t.Fatalf("%T does not implement EventWriter", r)
	}
}
