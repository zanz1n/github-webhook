package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var (
	configPath = flag.String(
		"config",
		"/var/github-webhook/config.yml",
		"The path of the configuration file",
	)
)

func init() {
	log.SetOutput(os.Stdout)
	flag.Parse()
}

func main() {
	enr, err := NewRepository(*configPath)
	if err != nil {
		log.Fatalln("Failed to load configuration:", err)
	}

	srv := NewServer(enr)

	if err = http.ListenAndServe(enr.GetAddr(), srv); err != nil {
		log.Fatalf(
			"Failed to listen on '%s': %s\n",
			enr.GetAddr(),
			err.Error(),
		)
	}
}
