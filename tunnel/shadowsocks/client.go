package shadowsocks

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/shadowsocks/go-shadowsocks2/core"

	"github.com/thep0y/trojan-go/common"
	"github.com/thep0y/trojan-go/config"
	"github.com/thep0y/trojan-go/tunnel"
)

type Client struct {
	underlay tunnel.Client
	core.Cipher
}

func (c *Client) DialConn(address *tunnel.Address, tunnel tunnel.Tunnel) (tunnel.Conn, error) {
	conn, err := c.underlay.DialConn(address, &Tunnel{})
	if err != nil {
		return nil, err
	}
	return &Conn{
		aeadConn: c.Cipher.StreamConn(conn),
		Conn:     conn,
	}, nil
}

func (c *Client) DialPacket(tunnel tunnel.Tunnel) (tunnel.PacketConn, error) {
	panic("not supported")
}

func (c *Client) Close() error {
	return c.underlay.Close()
}

func NewClient(ctx context.Context, underlay tunnel.Client) (*Client, error) {
	cfg := config.FromContext(ctx, Name).(*Config)
	cipher, err := core.PickCipher(cfg.Shadowsocks.Method, nil, cfg.Shadowsocks.Password)
	if err != nil {
		return nil, common.NewError("invalid shadowsocks cipher").Base(err)
	}
	log.Debug().Msg("shadowsocks client created")
	return &Client{
		underlay: underlay,
		Cipher:   cipher,
	}, nil
}
