package sse

import (
	"bytes"
	"errors"
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
	eventStream chan Event
	Closed      chan bool
}

func (es *EventStream) CloseNotify() <-chan bool {
	return es.writer.CloseNotify()
}

func (es *EventStream) writeHeaders() {
	es.writer.Header().Set("Content-Type", "text/event-stream")
	es.writer.Header().Set("Cache-Control", "no-cache")
	es.writer.Header().Set("Connection", "keep-alive")
	es.writer.Header().Set("X-Accel-Buffering", "no")
}

func (es *EventStream) writeEvent(event Event) {
	es.writer.Write(EventToBytes(event))
	es.writer.Flush()
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
			case <-es.writer.CloseNotify():
				es.Closed <- true
				return
			}
		}
	}()
}

func NewEventStream(w http.ResponseWriter) (*EventStream, error) {
	ew, ok := w.(EventWriter)
	if !ok {
		return nil, errors.New(fmt.Sprintf("%T doesn't supported streaming", w))
	}
	es := EventStream{
		writer:      ew,
		eventStream: make(chan Event),
	}
	return &es, nil
}
