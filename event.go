package sse

type Event interface {
	ID() string
	Event() string
	Data() string
}

type BasicEvent struct {
	id    string
	event string
	data  string
}

func (e BasicEvent) ID() string    { return e.id }
func (e BasicEvent) Event() string { return e.event }
func (e BasicEvent) Data() string  { return e.data }
