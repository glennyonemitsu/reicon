package server

import (
	"net"
	"path"
)

func CreateServer(basepath string) *Server {
	s := new(Server)
	s.Bootstrap(basepath)
	return s
}

func GetUnixAddress(basepath string) *net.UnixAddr {
	filepath := path.Join(basepath, "server.sock")
	addr := new(net.UnixAddr)
	addr.Name = filepath
	addr.Net = "unix"
	return addr
}
