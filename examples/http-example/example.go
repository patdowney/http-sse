package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	b := sse.StartNewBroker()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			b.Stop()
			// sig is a ^C, handle it
			time.Sleep(100 * time.Millisecond)
			os.Exit(0)
		}
	}()

	// generate some events
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
	/*
		// shutdown the broker if told to
		go func() {
			select {
			case <-time.After(30 * time.Second):
				b.Stop()
			}
		}()

		// start the broker again  for kicks
		go func() {
			select {
			case <-time.After(35 * time.Second):
				b.Start()
			}
		}()
	*/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s, err := sse.StartNewEventStream(w)
		if err != nil {
			log.Printf("error-starting-stream: %s", err.Error())
			return
		}
		b.Subscribe(s)

		for {
			select {
			case <-s.CloseNotify():
				b.Unsubscribe(s)
				return
			}
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
