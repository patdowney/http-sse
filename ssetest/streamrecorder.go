package ssetest

import (
	"net/http"
	"net/http/httptest"
)

type StreamRecorder struct {
	Recorder *httptest.ResponseRecorder
	closer   chan bool
}

// implement http.Flusher
func (w *StreamRecorder) Flush() {
	w.Recorder.Flush()
}

// implement http.CloseNotifier
func (w *StreamRecorder) CloseNotify() <-chan bool {
	return w.closer
}

// implement http.ResponseWriter
func (w *StreamRecorder) Header() http.Header {
	return w.Recorder.Header()
}

func (w *StreamRecorder) WriteHeader(code int) {
	w.Recorder.WriteHeader(code)
}

func (w *StreamRecorder) Write(data []byte) (int, error) {
	return w.Recorder.Write(data)
}

//Utility function to simulate client closing connection
func (w *StreamRecorder) Close() {
	w.closer <- true
}

func NewStreamRecorder() *StreamRecorder {
	r := StreamRecorder{
		Recorder: httptest.NewRecorder(),
		closer:   make(chan bool),
	}
	return &r
}
