package sse

import (
	"net/http"
)

type EventWriter interface {
	http.ResponseWriter
	http.Flusher
	//	http.CloseNotifier
}
