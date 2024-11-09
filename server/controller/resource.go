package controller

import (
	"github.com/kirilngusi/go-reverse-proxy/pkg/nathole"
	plugin "github.com/kirilngusi/go-reverse-proxy/pkg/plugin/server"
	"github.com/kirilngusi/go-reverse-proxy/pkg/util/tcpmux"
	"github.com/kirilngusi/go-reverse-proxy/pkg/util/vhost"
	"github.com/kirilngusi/go-reverse-proxy/server/group"
	"github.com/kirilngusi/go-reverse-proxy/server/ports"
	"github.com/kirilngusi/go-reverse-proxy/server/visitor"
)

// All resource managers and controllers
type ResourceController struct {
	// Manage all visitor listeners
	VisitorManager *visitor.Manager

	// TCP Group Controller
	TCPGroupCtl *group.TCPGroupCtl

	// HTTP Group Controller
	HTTPGroupCtl *group.HTTPGroupController

	// TCP Mux Group Controller
	TCPMuxGroupCtl *group.TCPMuxGroupCtl

	// Manage all TCP ports
	TCPPortManager *ports.Manager

	// Manage all UDP ports
	UDPPortManager *ports.Manager

	// For HTTP proxies, forwarding HTTP requests
	HTTPReverseProxy *vhost.HTTPReverseProxy

	// For HTTPS proxies, route requests to different clients by hostname and other information
	VhostHTTPSMuxer *vhost.HTTPSMuxer

	// Controller for nat hole connections
	NatHoleController *nathole.Controller

	// TCPMux HTTP CONNECT multiplexer
	TCPMuxHTTPConnectMuxer *tcpmux.HTTPConnectTCPMuxer

	// All server manager plugin
	PluginManager *plugin.Manager
}
