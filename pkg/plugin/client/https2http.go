//go:build !frps

package plugin

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	stdlog "log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/fatedier/golib/pool"
	"github.com/samber/lo"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/transport"
	httppkg "github.com/fatedier/frp/pkg/util/http"
	"github.com/fatedier/frp/pkg/util/log"
	netpkg "github.com/fatedier/frp/pkg/util/net"
)

func init() {
	Register(v1.PluginHTTPS2HTTP, NewHTTPS2HTTPPlugin)
}

type HTTPS2HTTPPlugin struct {
	opts *v1.HTTPS2HTTPPluginOptions

	l *Listener
	s *http.Server
}

func NewHTTPS2HTTPPlugin(options v1.ClientPluginOptions) (Plugin, error) {
	opts := options.(*v1.HTTPS2HTTPPluginOptions)
	listener := NewProxyListener()

	p := &HTTPS2HTTPPlugin{
		opts: opts,
		l:    listener,
	}

	rp := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.Out.Header["X-Forwarded-For"] = r.In.Header["X-Forwarded-For"]
			r.SetXForwarded()
			req := r.Out
			req.URL.Scheme = "http"
			req.URL.Host = p.opts.LocalAddr
			if p.opts.HostHeaderRewrite != "" {
				req.Host = p.opts.HostHeaderRewrite
			}
			for k, v := range p.opts.RequestHeaders.Set {
				req.Header.Set(k, v)
			}
		},
		BufferPool: pool.NewBuffer(32 * 1024),
		ErrorLog:   stdlog.New(log.NewWriteLogger(log.WarnLevel, 2), "", 0),
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil {
			tlsServerName, _ := httppkg.CanonicalHost(r.TLS.ServerName)
			host, _ := httppkg.CanonicalHost(r.Host)
			if tlsServerName != "" && tlsServerName != host {
				w.WriteHeader(http.StatusMisdirectedRequest)
				return
			}
		}
		rp.ServeHTTP(w, r)
	})

	tlsConfig, err := transport.NewServerTLSConfig(p.opts.CrtPath, p.opts.KeyPath, "")
	if err != nil {
		return nil, fmt.Errorf("gen TLS config error: %v", err)
	}

	p.s = &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 60 * time.Second,
		TLSConfig:         tlsConfig,
	}
	if !lo.FromPtr(opts.EnableHTTP2) {
		p.s.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
	}

	go func() {
		_ = p.s.ServeTLS(listener, "", "")
	}()
	return p, nil
}

func (p *HTTPS2HTTPPlugin) Handle(_ context.Context, conn io.ReadWriteCloser, realConn net.Conn, extra *ExtraInfo) {
	wrapConn := netpkg.WrapReadWriteCloserToConn(conn, realConn)
	if extra.SrcAddr != nil {
		wrapConn.SetRemoteAddr(extra.SrcAddr)
	}
	_ = p.l.PutConn(wrapConn)
}

func (p *HTTPS2HTTPPlugin) Name() string {
	return v1.PluginHTTPS2HTTP
}

func (p *HTTPS2HTTPPlugin) Close() error {
	return p.s.Close()
}
