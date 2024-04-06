package transport

import (
	"net"

	"github.com/thep0y/trojan-go/tunnel"
)

type Conn struct {
	net.Conn
}

func (c *Conn) Metadata() *tunnel.Metadata {
	return nil
}
