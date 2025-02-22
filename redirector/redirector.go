package redirector

import (
	"context"
	"io"
	"net"
	"reflect"

	"github.com/rs/zerolog/log"
)

type Dial func(net.Addr) (net.Conn, error)

func defaultDial(addr net.Addr) (net.Conn, error) {
	return net.Dial("tcp", addr.String())
}

type Redirection struct {
	Dial
	RedirectTo  net.Addr
	InboundConn net.Conn
}

type Redirector struct {
	ctx             context.Context
	redirectionChan chan *Redirection
}

func (r *Redirector) Redirect(redirection *Redirection) {
	select {
	case r.redirectionChan <- redirection:
		log.Debug().Msg("redirect request")
	case <-r.ctx.Done():
		log.Debug().Msg("exiting")
	}
}

func (r *Redirector) worker() {
	for {
		select {
		case redirection := <-r.redirectionChan:
			handle := func(redirection *Redirection) {
				if redirection.InboundConn == nil ||
					reflect.ValueOf(redirection.InboundConn).IsNil() {
					log.Error().Msg("nil inbound conn")
					return
				}
				defer redirection.InboundConn.Close()
				if redirection.RedirectTo == nil ||
					reflect.ValueOf(redirection.RedirectTo).IsNil() {
					log.Error().Msg("nil redirection addr")
					return
				}
				if redirection.Dial == nil {
					redirection.Dial = defaultDial
				}
				log.Warn().
					Stringer("from", redirection.InboundConn.RemoteAddr()).
					Stringer("to", redirection.RedirectTo).
					Msg("redirecting connection")
				outboundConn, err := redirection.Dial(redirection.RedirectTo)
				if err != nil {
					log.Error().Err(err).Msg("failed to redirect to target address")
					return
				}
				defer outboundConn.Close()
				errChan := make(chan error, 2)
				copyConn := func(a, b net.Conn) {
					_, err := io.Copy(a, b)
					errChan <- err
				}
				go copyConn(outboundConn, redirection.InboundConn)
				go copyConn(redirection.InboundConn, outboundConn)
				select {
				case err := <-errChan:
					if err != nil {
						log.Error().Err(err).Msg("failed to redirect")
					} else {
						log.Info().Msg("redirection done")
					}
				case <-r.ctx.Done():
					log.Debug().Msg("exiting")
					return
				}
			}
			go handle(redirection)
		case <-r.ctx.Done():
			log.Debug().Msg("shutting down redirector")
			return
		}
	}
}

func NewRedirector(ctx context.Context) *Redirector {
	r := &Redirector{
		ctx:             ctx,
		redirectionChan: make(chan *Redirection, 64),
	}
	go r.worker()
	return r
}
