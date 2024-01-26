package main

import (
	"io"
	"log/slog"
	"os"
	"os/exec"
	"sync"
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
	success := 0

	var (
		stdout io.ReadCloser
		buf    []byte
		err    error
	)
	for _, cmd := range cmds {
		stdout, err = exec.Command(cmd[0], cmd[1:]...).StdoutPipe()
		if err != nil {
			slog.Error(
				"Failed to execute action",
				"endpoint",
				endpoint,
				"event",
				event,
				"error",
				err,
			)

			continue
		}

		_, _ = os.Stderr.Write([]byte{'\n'})

		for {
			buf = make([]byte, 0, 512)

			_, err = stdout.Read(buf)
			if err != nil {
				if err == io.EOF {
					_, _ = os.Stderr.Write(buf)
					break
				} else {
					slog.Error(
						"Failed to execute action",
						"endpoint",
						endpoint,
						"event",
						event,
						"error",
						err,
					)

					break
				}
			}

			_, _ = os.Stderr.Write(buf)
		}

		_, _ = os.Stderr.Write([]byte{'\n'})
		_ = stdout.Close()
	}

	slog.Info(
		"Finished handling endpoint actions",
		"endpoint",
		endpoint,
		"event",
		event,
		"total_count",
		len(cmds),
		"success_count",
		success,
		"failed_count",
		len(cmds)-success,
	)
}
