package sse

import (
	"bytes"
	"context"

	"github.com/patdowney/http-sse/ssetest"
)

type TestStream struct {
	eventStream    *EventStream
	streamRecorder *ssetest.StreamRecorder
	cancelFunc     context.CancelFunc
}

func NewTestStream(start bool) (*TestStream, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	t := TestStream{
		streamRecorder: ssetest.NewStreamRecorder(),
		cancelFunc:     cancelFunc,
	}

	eventStream, err := NewEventStream(t.streamRecorder, ctx)
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
	t.cancelFunc()
	//	t.streamRecorder.Close()
}

func (t *TestStream) EventStream() *EventStream {
	return t.eventStream
}
