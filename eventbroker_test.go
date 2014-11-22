package sse

import (
	"runtime"
	"testing"
)

func TestBroker(t *testing.T) {
	b := NewBroker()
	b.Start()
	streamOne, errOne := NewTestStream(true)
	if errOne != nil {
		t.Fatalf("Failed to instantiate test stream one")
	}

	streamTwo, errTwo := NewTestStream(true)
	if errTwo != nil {
		t.Fatalf("Failed to instantiate test stream two")
	}

	b.Subscribe(streamOne.EventStream())
	b.Subscribe(streamTwo.EventStream())

	e1 := BasicEvent{id: "1", data: "dataone"}
	b.SendEvent(e1)

	b.Unsubscribe(streamOne.EventStream())

	e2 := BasicEvent{id: "2", data: "datatwo"}
	b.SendEvent(e2)

	b.Stop()

	// I am upset that I have resorted to this
	// need to let the other go routines run
	// so that the events can be received by
	// the streams
	runtime.Gosched()

	if !streamOne.Received(e1) {
		t.Errorf("streamOne(%p) failed to receive e1", streamOne.EventStream())
	}

	if !streamTwo.Received(e1) {
		t.Errorf("streamTwo(%p) failed to receive e1", streamTwo.EventStream())
	}

	if !streamTwo.Received(e2) {
		t.Errorf("streamTwo(%p) failed to receive e2", streamTwo.EventStream())
	}

}
