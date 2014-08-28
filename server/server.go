package server

import (
	"encoding/gob"
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

	s.signalChan = make(chan os.Signal)
	signal.Notify(s.signalChan)

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
	s.SocketAddress = GetUnixAddress(s.basepath)
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
	go s.ListenSignal()
	s.ListenUnixConnection()
	return nil
}

func (s *Server) handleUnixConn(conn *net.UnixConn) {
	var args []string
	dec := gob.NewDecoder(conn)
	dec.Decode(&args)

	for _, a := range args {
		println(a)
	}
	conn.Close()
}

func (s *Server) ListenUnixConnection() {
	for {
		conn, err := s.Socket.AcceptUnix()
		if err != nil {
			continue
		}
		go s.handleUnixConn(conn)
	}
}
func (s *Server) ListenSignal() {
	var si os.Signal
	for {
		si = <-s.signalChan
		switch si {
		case syscall.SIGINT:
			fallthrough
		case syscall.SIGQUIT:
			s.Shutdown()
		}
	}
}

func (s *Server) Shutdown() {
	s.CloseSocket()
	os.Exit(0)
}
