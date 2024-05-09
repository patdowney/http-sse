package sse

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/patdowney/http-sse/ssetest"
)

func TestEventToBytes(t *testing.T) {
	e1 := BasicEvent{id: "1", data: "dataone"}
	e1ExpectedBytes := []byte("id: 1\ndata: dataone\n\n")

	e1Bytes := EventToBytes(e1)
	if !bytes.Equal(e1Bytes, e1ExpectedBytes) {
		t.Fatalf("unexpected serialisation:\nexpected:\n>\n%s\n<\nreceived:\n>\n%s\n<\n", e1ExpectedBytes, e1Bytes)
	}

	e2 := BasicEvent{id: "2", event: "eventtwo", data: "datatwo"}
	e2ExpectedBytes := []byte("id: 2\nevent: eventtwo\ndata: datatwo\n\n")

	e2Bytes := EventToBytes(e2)
	if !bytes.Equal(e2Bytes, e2ExpectedBytes) {
		t.Fatalf("unexpected serialisation:\nexpected:\n>\n%s\n<\nreceived:\n>\n%s\n<\n", e2ExpectedBytes, e2Bytes)
	}
}

func TestEventStream(t *testing.T) {
	w := ssetest.NewStreamRecorder()

	ctx, cancel := context.WithCancel(context.Background())

	s, err := NewEventStream(w, ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	s.Start()

	s.SendEvent(BasicEvent{id: "1", data: "dataone"})
	s.SendEvent(BasicEvent{id: "2", event: "omg", data: "datatwo"})

	time.Sleep(5 * time.Millisecond)
	cancel()

	if w.Recorder.Code != http.StatusOK {
		t.Error("Expected 200 response")
	}
	if !w.Recorder.Flushed {
		t.Error("Response not flushed")
	}

	expectedOutput := "id: 1\ndata: dataone\n\nid: 2\nevent: omg\ndata: datatwo\n\n"

	if expectedOutput != string(w.Bytes()) {
		t.Errorf("expected:\n%s\nreceived:\n%s\n", expectedOutput, w.Recorder.Body)
	}
}
