package app

import (
	"github.com/net-byte/vtun/dtls"
	"github.com/net-byte/vtun/h1"
	"github.com/net-byte/vtun/h2"
	"github.com/net-byte/vtun/kcp"
	"github.com/net-byte/vtun/quic"
	"github.com/net-byte/vtun/tcp"
	"github.com/net-byte/vtun/utls"
	"log"

	"github.com/net-byte/vtun/common/cipher"
	"github.com/net-byte/vtun/common/config"
	"github.com/net-byte/vtun/common/netutil"
	"github.com/net-byte/vtun/grpc"
	"github.com/net-byte/vtun/tls"
	"github.com/net-byte/vtun/tun"
	"github.com/net-byte/vtun/udp"
	"github.com/net-byte/vtun/ws"
	"github.com/net-byte/water"
)

var _banner = `
_                 
__ __ | |_   _  _   _ _  
\ V / |  _| | || | | ' \ 
 \_/   \__|  \_,_| |_||_|
						 
A simple VPN written in Go.
%s
`
var _srcUrl = "https://github.com/net-byte/vtun"

// vtun app struct
type App struct {
	Config  *config.Config
	Version string
	Iface   *water.Interface
}

func NewApp(config *config.Config, version string) *App {

	return &App{
		Config:  config,
		Version: version,
	}
}

// InitConfig initializes the config
func (app *App) InitConfig() {
	log.Printf(_banner, _srcUrl)
	log.Printf("vtun version %s", app.Version)
	if !app.Config.ServerMode {
		app.Config.LocalGateway = netutil.DiscoverGateway(true)
		app.Config.LocalGatewayv6 = netutil.DiscoverGateway(false)
	}
	app.Config.BufferSize = 64 * 1024
	cipher.SetKey(app.Config.Key)
	app.Iface = tun.CreateTun(*app.Config)
	log.Printf("initialized config: %+v", app.Config)
	netutil.PrintStats(app.Config.Verbose, app.Config.ServerMode)
}

// StartApp starts the app
func (app *App) StartApp() {

	switch app.Config.Protocol {
	case "udp":
		if app.Config.ServerMode {
			udp.StartServer(app.Iface, *app.Config)
		} else {
			udp.StartClient(app.Iface, *app.Config)
		}
	case "ws", "wss":
		if app.Config.ServerMode {
			ws.StartServer(app.Iface, *app.Config)
		} else {
			ws.StartClient(app.Iface, *app.Config)
		}
	case "tls":
		if app.Config.ServerMode {
			tls.StartServer(app.Iface, *app.Config)
		} else {
			tls.StartClient(app.Iface, *app.Config)
		}
	case "grpc":
		if app.Config.ServerMode {
			grpc.StartServer(app.Iface, *app.Config)
		} else {
			grpc.StartClient(app.Iface, *app.Config)
		}
	case "quic":
		if app.Config.ServerMode {
			quic.StartServer(app.Iface, *app.Config)
		} else {
			quic.StartClient(app.Iface, *app.Config)
		}
	case "kcp":
		if app.Config.ServerMode {
			kcp.StartServer(app.Iface, *app.Config)
		} else {
			kcp.StartClient(app.Iface, *app.Config)
		}
	case "utls":
		if app.Config.ServerMode {
			utls.StartServer(app.Iface, *app.Config)
		} else {
			utls.StartClient(app.Iface, *app.Config)
		}
	case "dtls":
		if app.Config.ServerMode {
			dtls.StartServer(app.Iface, *app.Config)
		} else {
			dtls.StartClient(app.Iface, *app.Config)
		}
	case "h2":
		if app.Config.ServerMode {
			h2.StartServer(app.Iface, *app.Config)
		} else {
			h2.StartClient(app.Iface, *app.Config)
		}
	case "tcp":
		if app.Config.ServerMode {
			tcp.StartServer(app.Iface, *app.Config)
		} else {
			tcp.StartClient(app.Iface, *app.Config)
		}
	case "http":
		if app.Config.ServerMode {
			h1.StartServer(app.Iface, *app.Config)
		} else {
			h1.StartClient(app.Iface, *app.Config)
		}
	default:
		if app.Config.ServerMode {
			udp.StartServer(app.Iface, *app.Config)
		} else {
			udp.StartClient(app.Iface, *app.Config)
		}
	}
}

// StopApp stops the app
func (app *App) StopApp() {
	tun.ResetRoute(*app.Config)
	app.Iface.Close()
	log.Println("vtun stopped")
}
