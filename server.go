package main

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/webhooks/v6/github"
)

func NewServer(er *Repository) *Server {
	return &Server{er}
}

type Server struct {
	er *Repository
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	endpoint, ok := s.er.GetEndpoint(r.URL.Path)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	signature := r.Header.Get("X-Hub-Signature")
	if signature == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event := r.Header.Get("X-GitHub-Event")
	if event == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slog.Info("Incomming webhook request", "endpoint", endpoint.Route, "event", event)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ValidatePayload(endpoint.Secret, payload, signature) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cmds := endpoint.GetCommands(github.Event(event))
	if len(cmds) > 0 {
		slog.Info(
			"Triggering endpoint actions",
			"endpoint",
			endpoint.Route,
			"event",
			event,
			"actions_count",
			len(cmds),
		)
		go s.er.HandleCommands(cmds, endpoint.Route, event)
	} else {
		slog.Info(
			"Endpoint has no actions setted up for the event",
			"endpoint",
			endpoint.Route,
			"event",
			event,
		)
	}

	w.WriteHeader(http.StatusOK)
}
