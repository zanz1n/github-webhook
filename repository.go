package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"time"
)

func NewRepository(configPath string) (*Repository, error) {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &Repository{
		config:     cfg,
		configMu:   sync.RWMutex{},
		configPath: configPath,
	}, nil
}

type Repository struct {
	config     *Config
	configMu   sync.RWMutex
	configPath string
}

func (r *Repository) RefreshConfig() error {
	cfg, err := LoadConfig(r.configPath)
	if err != nil {
		return err
	}

	r.configMu.Lock()
	*r.config = *cfg
	r.configMu.Unlock()

	return nil
}

func (r *Repository) GetAddr() string {
	r.configMu.RLock()
	s := r.config.Addr
	r.configMu.RUnlock()

	return s
}

func (r *Repository) GetEndpoint(route string) (*Endpoint, bool) {
	r.configMu.RLock()
	defer r.configMu.RUnlock()

	for _, e := range r.config.Endpoints {
		if e.Route == route {
			return &e, true
		}
	}
	return nil, false
}

func (r *Repository) HandleCommands(cmds [][]string, endpoint, event string) {
	timeStart := time.Now()

	success := 0
	l := len(cmds)

	var (
		out         bytes.Buffer
		actionStart time.Time
	)
	for i, cmd := range cmds {
		actionStart = time.Now()

		action := fmt.Sprintf("%d/%d", i+1, l)
		cmd := exec.Command(cmd[0], cmd[1:]...)

		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			slog.Error(
				"Failed to execute action",
				"endpoint",
				endpoint,
				"event",
				event,
				"action",
				action,
				"error",
				err,
			)
		} else {
			_, _ = os.Stderr.Write([]byte{'\n'})
			_, _ = os.Stderr.Write(out.Bytes())
			_, _ = os.Stderr.Write([]byte{'\n'})

			slog.Info(
				"Action executed successfully",
				"endpoint",
				endpoint,
				"event",
				event,
				"action",
				action,
				"time_taken",
				time.Since(actionStart),
			)
			success++
		}

		out.Reset()
	}

	slog.Info(
		"Finished handling actions",
		"endpoint",
		endpoint,
		"event",
		event,
		"total_count",
		l,
		"success_count",
		success,
		"failed_count",
		l-success,
		"time_taken",
		time.Since(timeStart),
	)
}
