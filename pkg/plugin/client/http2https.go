package plugin

import (
	"context"
	"crypto/tls"
	"io"
	stdlog "log"
	"net"
	"net/http"
	"net/http/httputil"

	"github.com/fatedier/golib/pool"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/util/log"
	netpkg "github.com/fatedier/frp/pkg/util/net"
)

func init() {
	Register(v1.PluginHTTP2HTTPS, NewHTTP2HTTPSPlugin)
}

type HTTP2HTTPSPlugin struct {
	opts *v1.HTTP2HTTPSPluginOptions

	l *Listener
	s *http.Server
}

func NewHTTP2HTTPSPlugin(options v1.ClientPluginOptions) (Plugin, error) {
	opts := options.(*v1.HTTP2HTTPSPluginOptions)

	listener := NewProxyListener()

	p := &HTTP2HTTPSPlugin{
		opts: opts,
		l:    listener,
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	rp := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.Out.Header["X-Forwarded-For"] = r.In.Header["X-Forwarded-For"]
			r.Out.Header["X-Forwarded-Host"] = r.In.Header["X-Forwarded-Host"]
			r.Out.Header["X-Forwarded-Proto"] = r.In.Header["X-Forwarded-Proto"]
			req := r.Out
			req.URL.Scheme = "https"
			req.URL.Host = p.opts.LocalAddr
			if p.opts.HostHeaderRewrite != "" {
				req.Host = p.opts.HostHeaderRewrite
			}
			for k, v := range p.opts.RequestHeaders.Set {
				req.Header.Set(k, v)
			}
		},
		Transport:  tr,
		BufferPool: pool.NewBuffer(32 * 1024),
		ErrorLog:   stdlog.New(log.NewWriteLogger(log.WarnLevel, 2), "", 0),
	}

	p.s = &http.Server{
		Handler:           rp,
		ReadHeaderTimeout: 0,
	}

	go func() {
		_ = p.s.Serve(listener)
	}()

	return p, nil
}

func (p *HTTP2HTTPSPlugin) Handle(_ context.Context, conn io.ReadWriteCloser, realConn net.Conn, _ *ExtraInfo) {
	wrapConn := netpkg.WrapReadWriteCloserToConn(conn, realConn)
	_ = p.l.PutConn(wrapConn)
}

func (p *HTTP2HTTPSPlugin) Name() string {
	return v1.PluginHTTP2HTTPS
}

func (p *HTTP2HTTPSPlugin) Close() error {
	return p.s.Close()
}
