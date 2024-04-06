package server

import (
	"github.com/thep0y/trojan-go/config"
	"github.com/thep0y/trojan-go/proxy/client"
)

func init() {
	config.RegisterConfigCreator(Name, func() interface{} {
		return new(client.Config)
	})
}
