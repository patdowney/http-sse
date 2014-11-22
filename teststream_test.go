package sse

import (
	"testing"
)

func TestTestStream(t *testing.T) {
	ts, err := NewTestStream(true)
	if err != nil {
		t.Fatalf("Failed to instantiate TestStream")
	}

	e1 := BasicEvent{id: "1", data: "dataone"}
	e2 := BasicEvent{id: "2", data: "datatwo"}

	ts.EventStream().SendEvent(e1)
	ts.EventStream().SendEvent(e2)

	ts.CloseRecorder()

	if !ts.Received(e1) {
		t.Errorf("Didn't receive e1")
	}

	if !ts.Received(e2) {
		t.Errorf("Didn't receive e2")
	}
}
