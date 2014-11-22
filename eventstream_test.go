package sse

import (
	//	"fmt"
	"net/http"
	"testing"

	"github.com/patdowney/http-sse/ssetest"
)

func TestEventStream(t *testing.T) {
	w := ssetest.NewStreamRecorder()

	s, err := NewEventStream(w)
	if err != nil {
		t.Fatalf(err.Error())
	}

	s.Start()

	s.SendEvent(BasicEvent{id: "1", data: "dataone"})
	s.SendEvent(BasicEvent{id: "2", event: "omg", data: "datatwo"})

	w.Close()

	if w.Recorder.Code != http.StatusOK {
		t.Error("Expected 200 response")
	}
	if !w.Recorder.Flushed {
		t.Error("Response not flushed")
	}

	expectedOutput := "id: 1\ndata: dataone\n\nid: 2\nevent: omg\ndata: datatwo\n\n"

	if expectedOutput != w.Recorder.Body.String() {
		t.Errorf("expected:\n%s\nreceived:\n%s\n", expectedOutput, w.Recorder.Body)
	}
}
