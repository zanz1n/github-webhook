package main

import "github.com/go-playground/webhooks/v6/github"

type Config struct {
	Addr      string     `json:"addr" yaml:"addr" validate:"required"`
	Endpoints []Endpoint `json:"endpoints" yaml:"endpoints" validate:"required"`
}

type Endpoint struct {
	Route   string   `json:"route" yaml:"route" validate:"required,dirpath"`
	Secret  string   `json:"secret" yaml:"secret" validate:"required"`
	Actions []Action `json:"actions" yaml:"actions" validate:"required"`
}

type Action struct {
	Event   github.Event `json:"event" yaml:"event" validate:"required"`
	Command []string     `json:"command" yaml:"command" validate:"required"`
}

func (e *Endpoint) GetCommands(event github.Event) [][]string {
	cmds := [][]string{}

	for _, a := range e.Actions {
		if a.Event == event {
			cmds = append(cmds, a.Command)
		}
	}

	return cmds
}
