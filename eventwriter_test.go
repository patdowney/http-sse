package sse

import (
	"testing"

	"github.com/patdowney/http-sse/ssetest"
)

func TestEnsureStreamRecorderImplementsEventWriter(t *testing.T) {
	r := ssetest.NewStreamRecorder()
	i := interface{}(r)
	_, ok := i.(EventWriter)
	if !ok {
		t.Fatalf("%T does not implement EventWriter", r)
	}
}
