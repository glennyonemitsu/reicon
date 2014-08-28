package server

import (
	"encoding/gob"
	"fmt"
	"log"
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

	err = s.OpenSocket()
	if err != nil {
		return err
	}

	err = s.FindModules()
	if err != nil {
		return err
	}

	err = s.LoadModules()
	if err != nil {
		return err
	}
	go s.PingModules()

	return nil
}

func (s *Server) PingModules() {
	log.Println("pinging")

	conn, err := net.DialUnix("unix", nil, s.SocketAddress)
	if err != nil {
		log.Fatalf("Error when dialing for ping: %s\n", err.Error())
	}

	enc := gob.NewEncoder(conn)
	err = enc.Encode(string("ping"))
	if err != nil {
		log.Fatalf("could not encode %s\n", err.Error())
	}

	var reply string
	dec := gob.NewDecoder(conn)
	err = dec.Decode(&reply)
	if err != nil {
		log.Fatalf("could not decode %s\n", err.Error())
	}
	log.Println(reply)
}

func (s *Server) FindModules() error {
	log.Println("finding modules")
	var err error
	s.Modules, err = filepath.Glob(path.Join(s.modpath, "*"))
	if err != nil {
		return err
	}
	log.Println("found modules")
	return nil
}

func (s *Server) LoadModules() error {
	log.Println("loading modules")
	var proc *ModuleProc
	var attr *os.ProcAttr
	var err error
	for _, m := range s.Modules {
		proc = new(ModuleProc)
		attr = new(os.ProcAttr)
		attr.Dir = m
		attr.Env = append(
			attr.Env,
			fmt.Sprintf("REICON_SYSTEM_SOCKET=%s", s.socketpath),
		)
		proc.Process, err = os.StartProcess(
			fmt.Sprintf("%s/mod_%s", m, path.Base(m)),
			[]string{},
			attr,
		)
		if err != nil {
			return err
		}
		s.ModuleProcs = append(s.ModuleProcs, proc)
	}
	log.Println("loaded modules")
	return nil
}

func (s *Server) OpenSocket() error {
	s.SocketAddress = GetUnixAddress(s.basepath)
	listen, err := net.ListenUnix("unix", s.SocketAddress)
	if err != nil {
		return err
	}
	s.SocketListen = listen
	return nil
}

func (s *Server) CloseSocket() error {
	if s.SocketListen != nil {
		return s.SocketListen.Close()
	}
	return nil
}

func (s *Server) Run() {
	go s.ListenSignal()
	s.ListenUnixConnection()
}

func (s *Server) handleUnixConn(conn *net.UnixConn) {
	var args []string
	dec := gob.NewDecoder(conn)
	dec.Decode(&args)

	for _, a := range args {
		log.Println(a)
	}
	conn.Close()
}

func (s *Server) ListenUnixConnection() {
	for {
		conn, err := s.SocketListen.AcceptUnix()
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
