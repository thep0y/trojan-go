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
	flag.Parse()

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
