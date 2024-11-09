package plugin

import (
	"context"
	"io"
	"net"

	libio "github.com/fatedier/golib/io"

	v1 "github.com/kirilngusi/go-reverse-proxy/pkg/config/v1"
	"github.com/kirilngusi/go-reverse-proxy/pkg/util/xlog"
)

func init() {
	Register(v1.PluginUnixDomainSocket, NewUnixDomainSocketPlugin)
}

type UnixDomainSocketPlugin struct {
	UnixAddr *net.UnixAddr
}

func NewUnixDomainSocketPlugin(options v1.ClientPluginOptions) (p Plugin, err error) {
	opts := options.(*v1.UnixDomainSocketPluginOptions)

	unixAddr, errRet := net.ResolveUnixAddr("unix", opts.UnixPath)
	if errRet != nil {
		err = errRet
		return
	}

	p = &UnixDomainSocketPlugin{
		UnixAddr: unixAddr,
	}
	return
}

func (uds *UnixDomainSocketPlugin) Handle(ctx context.Context, conn io.ReadWriteCloser, _ net.Conn, extra *ExtraInfo) {
	xl := xlog.FromContextSafe(ctx)
	localConn, err := net.DialUnix("unix", nil, uds.UnixAddr)
	if err != nil {
		xl.Warnf("dial to uds %s error: %v", uds.UnixAddr, err)
		return
	}
	if extra.ProxyProtocolHeader != nil {
		if _, err := extra.ProxyProtocolHeader.WriteTo(localConn); err != nil {
			return
		}
	}

	libio.Join(localConn, conn)
}

func (uds *UnixDomainSocketPlugin) Name() string {
	return v1.PluginUnixDomainSocket
}

func (uds *UnixDomainSocketPlugin) Close() error {
	return nil
}
