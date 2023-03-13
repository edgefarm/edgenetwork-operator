package main

import (
	golog "log"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"flag"

	sync "github.com/edgefarm/edgenetwork-operator/pkg/sync"
)

func main() {
	loglevel := flag.String("log-level", "info", "log level (trace, debug, info, warn, error, fatal, panic)")
	flag.Parse()
	l, err := zerolog.ParseLevel(*loglevel)
	if err != nil {
		golog.Fatal(err)
	}

	zerolog.SetGlobalLevel(l)
	log.Info().Msg("Starting hook")

	http.HandleFunc("/sync", sync.Handler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal().Msgf("Error starting server: %v", err)
	}
}
