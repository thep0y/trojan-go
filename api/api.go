package api

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/thep0y/trojan-go/statistic"
)

type Handler func(ctx context.Context, auth statistic.Authenticator) error

var handlers = make(map[string]Handler)

func RegisterHandler(name string, handler Handler) {
	handlers[name] = handler
}

func RunService(ctx context.Context, name string, auth statistic.Authenticator) error {
	if h, ok := handlers[name]; ok {
		log.Debug().Str("name", name).Msg("api handler found")
		return h(ctx, auth)
	}
	log.Debug().Str("name", name).Msg("api handler not found")
	return nil
}
