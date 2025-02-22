package main

import (
	"flag"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/thep0y/trojan-go/component"
	"github.com/thep0y/trojan-go/option"
)

const (
	TimeFormat = "2006-01-02 15:04:05.999999999"
)

func init() {
	zerolog.TimeFieldFormat = TimeFormat
}

func main() {
	log.Logger = log.With().Caller().Logger()

	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	for {
		h, err := option.PopOptionHandler()
		if err != nil {
			log.Fatal().Msg("invalid options")
		}
		err = h.Handle()
		if err == nil {
			break
		}
	}
}
