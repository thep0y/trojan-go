//go:build custom || full
// +build custom full

package build

import (
	_ "github.com/thep0y/trojan-go/proxy/custom"
	_ "github.com/thep0y/trojan-go/tunnel/adapter"
	_ "github.com/thep0y/trojan-go/tunnel/dokodemo"
	_ "github.com/thep0y/trojan-go/tunnel/freedom"
	_ "github.com/thep0y/trojan-go/tunnel/http"
	_ "github.com/thep0y/trojan-go/tunnel/mux"
	_ "github.com/thep0y/trojan-go/tunnel/router"
	_ "github.com/thep0y/trojan-go/tunnel/shadowsocks"
	_ "github.com/thep0y/trojan-go/tunnel/simplesocks"
	_ "github.com/thep0y/trojan-go/tunnel/socks"
	_ "github.com/thep0y/trojan-go/tunnel/tls"
	_ "github.com/thep0y/trojan-go/tunnel/tproxy"
	_ "github.com/thep0y/trojan-go/tunnel/transport"
	_ "github.com/thep0y/trojan-go/tunnel/trojan"
	_ "github.com/thep0y/trojan-go/tunnel/websocket"
)
