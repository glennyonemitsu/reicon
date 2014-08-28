package server

import (
	"net"
)

type Server struct {
	basepath      string
	modpath       string
	socketpath    string
	Modules       []string
	Socket        *net.UnixListener
	SocketAddress *net.UnixAddr
}
