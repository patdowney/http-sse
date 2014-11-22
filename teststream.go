package sse

import (
	"bytes"

	"github.com/patdowney/http-sse/ssetest"
)

type TestStream struct {
	eventStream    *EventStream
	streamRecorder *ssetest.StreamRecorder
}

func NewTestStream(start bool) (*TestStream, error) {
	t := TestStream{}
	t.streamRecorder = ssetest.NewStreamRecorder()

	eventStream, err := NewEventStream(t.streamRecorder)
	if err != nil {
		return nil, err
	}

	t.eventStream = eventStream

	if start {
		t.eventStream.Start()
	}

	return &t, nil
}

func (t *TestStream) Received(event Event) bool {
	eventBytes := EventToBytes(event)
	streamBytes := t.streamRecorder.Bytes()

	return bytes.Contains(streamBytes, eventBytes)
}

func (t *TestStream) CloseRecorder() {
	t.streamRecorder.Close()
}

func (t *TestStream) EventStream() *EventStream {
	return t.eventStream
}
