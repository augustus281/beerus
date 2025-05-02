package gateway

import (
	"net/http"
)

type WebsocketUpgrader func(w http.ResponseWriter, r *http.Request)

func NewServeMux(opts ...Option) http.Handler {
	mux := http.NewServeMux()
	var upgrader WebsocketUpgrader

	for _, opt := range opts {
		opt(&upgrader)
	}

	if upgrader != nil {
		mux.HandleFunc("/ws", upgrader)
	}

	// add gRPC-Gateway routes here
	return mux
}

type Option func(*WebsocketUpgrader)

func WithWebsocketUpgrader(u WebsocketUpgrader) Option {
	return func(target *WebsocketUpgrader) {
		*target = u
	}
}
