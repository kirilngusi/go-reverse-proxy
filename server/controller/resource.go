package controller

import (
	"github.com/fatedier/frp/pkg/nathole"
	plugin "github.com/fatedier/frp/pkg/plugin/server"
	"github.com/fatedier/frp/pkg/util/tcpmux"
	"github.com/fatedier/frp/pkg/util/vhost"
	"github.com/fatedier/frp/server/group"
	"github.com/fatedier/frp/server/ports"
	"github.com/fatedier/frp/server/visitor"
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
