package sse

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type EventStream struct {
	writer      EventWriter
	eventStream chan Event
	Closed      chan bool
}

func (es *EventStream) writeHeaders() {
	es.writer.Header().Set("Content-Type", "text/event-stream")
	es.writer.Header().Set("Cache-Control", "no-cache")
	es.writer.Header().Set("Connection", "keep-alive")
	es.writer.Header().Set("X-Accel-Buffering", "no")
}

func (es *EventStream) writeEvent(event Event) {
	if event.ID() != "" {
		fmt.Fprintf(es.writer, "id: %s\n", event.ID())
	}
	if event.Event() != "" {
		fmt.Fprintf(es.writer, "event: %s\n", event.Event())
	}
	fmt.Fprintf(es.writer, "data: %s\n\n", event.Data())
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
				log.Println("Closing connection")
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
