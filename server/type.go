package server

import (
	"net"
	"os"
)

type Server struct {
	basepath      string
	modpath       string
	socketpath    string
	signalChan    chan os.Signal
	connChan      chan *net.UnixConn
	Modules       []string
	Socket        *net.UnixListener
	SocketAddress *net.UnixAddr
}
