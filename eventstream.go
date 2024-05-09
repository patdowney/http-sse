package sse

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

func EventToBytes(event Event) []byte {
	var b bytes.Buffer

	if event.ID() != "" {
		fmt.Fprintf(&b, "id: %s\n", event.ID())
	}
	if event.Event() != "" {
		fmt.Fprintf(&b, "event: %s\n", event.Event())
	}
	fmt.Fprintf(&b, "data: %s\n\n", event.Data())

	return b.Bytes()
}

type EventStream struct {
	writer      EventWriter
	context     context.Context
	eventStream chan Event
	closed      chan bool
}

func (es *EventStream) writeHeaders() {
	es.writer.Header().Set("Content-Type", "text/event-stream")
	es.writer.Header().Set("Cache-Control", "no-cache")
	es.writer.Header().Set("Connection", "keep-alive")
	es.writer.Header().Set("X-Accel-Buffering", "no")
}

func (es *EventStream) writeEvent(event Event) {
	defer func() {
		es.writer.Flush()
	}()
	_, _ = es.writer.Write(EventToBytes(event))
}

func (es *EventStream) SendEvent(event Event) {
	es.eventStream <- event
}

func (es *EventStream) Start() {
	es.writeHeaders()

	go func() {
		for {
			select {
			case event := <-es.eventStream:
				es.writeEvent(event)
			case <-es.context.Done():
				es.closed <- true
				return
			}
		}
	}()
}

func (es *EventStream) Stop() {
	es.closed <- true
}

func NewEventStream(w http.ResponseWriter, c context.Context) (*EventStream, error) {
	ew, ok := w.(EventWriter)
	if !ok {
		return nil, fmt.Errorf("%T doesn't supported streaming", w)
	}
	es := EventStream{
		writer:      ew,
		context:     c,
		eventStream: make(chan Event),
		closed:      make(chan bool),
	}
	return &es, nil
}

func StartNewEventStream(w http.ResponseWriter, c context.Context) (*EventStream, error) {
	s, err := NewEventStream(w, c)
	if err != nil {
		return nil, err
	}
	s.Start()

	return s, nil
}
