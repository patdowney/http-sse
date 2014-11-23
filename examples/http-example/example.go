package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/patdowney/http-sse"
)

type ExampleEvent struct {
	id    string
	event string
	data  string
}

func (e ExampleEvent) ID() string    { return e.id }
func (e ExampleEvent) Event() string { return e.event }
func (e ExampleEvent) Data() string  { return e.data }

func main() {
	b := sse.NewBroker()
	b.Start()

	go func() {
		c := time.Tick(200 * time.Millisecond)
		id := 0
		for now := range c {
			data := fmt.Sprintf("%v", now)
			e := ExampleEvent{id: fmt.Sprintf("%v", id), data: data}
			b.SendEvent(e)
			id += 1
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s, _ := sse.NewEventStream(w)
		s.Start()
		b.Subscribe(s)

		select {
		case <-s.CloseNotify():
			b.Unsubscribe(s)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
