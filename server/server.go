package server

import (
	"net"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
)

func (s *Server) Bootstrap(basepath string) error {
	var err error
	s.basepath = basepath
	s.modpath = path.Join(basepath, "modules")
	s.socketpath = path.Join(basepath, "server.sock")
	err = s.FindModules()
	if err != nil {
		return err
	}
	err = s.OpenSocket()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) FindModules() error {
	var err error
	s.Modules, err = filepath.Glob(path.Join(s.modpath, "*"))
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) OpenSocket() error {
	s.SocketAddress = new(net.UnixAddr)
	s.SocketAddress.Name = s.socketpath
	s.SocketAddress.Net = "unix"
	socket, err := net.ListenUnix("unix", s.SocketAddress)
	if err != nil {
		return err
	}
	s.Socket = socket
	return nil
}

func (s *Server) CloseSocket() error {
	if s.Socket != nil {
		return s.Socket.Close()
	}
	return nil
}

func (s *Server) Run() error {
	var signalIn os.Signal
	var endSignal chan os.Signal
	endSignal = make(chan os.Signal)
	signal.Notify(endSignal, syscall.SIGINT, syscall.SIGQUIT)
	for {
		select {
		case signalIn = <-endSignal:
			println(signalIn.String())
			s.Shutdown()
			return nil
		}
	}
	return nil
}

func (s *Server) Shutdown() {
	s.CloseSocket()
	println("closing")
}
