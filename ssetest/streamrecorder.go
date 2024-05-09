package ssetest

import (
	"net/http"
	"net/http/httptest"
	"sync"
)

type StreamRecorder struct {
	Recorder  *httptest.ResponseRecorder
	closer    chan bool
	dataMutex *sync.Mutex
}

// implement http.Flusher
func (w *StreamRecorder) Flush() {
	w.Recorder.Flush()
}

// implement http.ResponseWriter
func (w *StreamRecorder) Header() http.Header {
	return w.Recorder.Header()
}

func (w *StreamRecorder) WriteHeader(code int) {
	w.Recorder.WriteHeader(code)
}

func (w *StreamRecorder) Write(data []byte) (int, error) {
	w.dataMutex.Lock()
	written, err := w.Recorder.Write(data)
	w.dataMutex.Unlock()

	return written, err
}

//Utility function to simulate client closing connection
func (w *StreamRecorder) Close() {
	w.closer <- true
}

func (w *StreamRecorder) Bytes() []byte {
	w.dataMutex.Lock()
	data := []byte(w.Recorder.Body.String())
	w.dataMutex.Unlock()

	return data
}

func NewStreamRecorder() *StreamRecorder {
	r := StreamRecorder{
		Recorder:  httptest.NewRecorder(),
		closer:    make(chan bool),
		dataMutex: &sync.Mutex{},
	}
	return &r
}
