package memory

import (
	"github.com/thep0y/trojan-go/config"
)

type Config struct {
	Passwords []string `json:"password" yaml:"password"`
}

func init() {
	config.RegisterConfigCreator(Name, func() interface{} {
		return &Config{}
	})
}
