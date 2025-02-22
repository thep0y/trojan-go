package redirector

import (
	"context"
	"io"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"

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
				var wg sync.WaitGroup
				copyConn := func(dst, src net.Conn) {
					defer wg.Done()
					_, err := io.Copy(dst, src)
					if err != nil {
						if err != io.EOF && !strings.Contains(err.Error(), "connection reset by peer") {
							log.Debug().
								Err(err).
								Str("from", src.RemoteAddr().String()).
								Str("to", dst.RemoteAddr().String()).
								Msg("connection copy error")
						}
					}
					// 确保连接被正确关闭
					dst.SetDeadline(time.Now())
					src.SetDeadline(time.Now())
					errChan <- err
				}

				wg.Add(2)
				go copyConn(outboundConn, redirection.InboundConn)
				go copyConn(redirection.InboundConn, outboundConn)

				go func() {
					wg.Wait()
					close(errChan)
				}()

				for err := range errChan {
					if err != nil && err != io.EOF && !strings.Contains(err.Error(), "connection reset by peer") {
						log.Error().Err(err).Msg("failed to redirect")
					}
				}
				log.Info().Msg("redirection done")
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
