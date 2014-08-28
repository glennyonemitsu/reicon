package module

import (
	"encoding/gob"
	"fmt"
	"github.com/glennyonemitsu/reicon/system/activity"
	"net"
	"os"
)

type Module struct {
	Name         string
	Activities   []*activity.Activity
	ServerSocket *net.UnixConn
	SocketAddr   *net.UnixAddr
	Listener     *net.UnixListener
}

func (m *Module) Run() {
	socket_path := os.Getenv("REICON_SYSTEM_SOCKET")
	m.SocketAddr = new(net.UnixAddr)
	m.SocketAddr.Name = "unix"
	m.SocketAddr.Net = socket_path
	m.Listener, _ = net.ListenUnix("unix", m.SocketAddr)
	m.ListenSocket()

}

func (m *Module) ListenSocket() {
	for {
		conn, err := m.Listener.AcceptUnix()
		if err != nil {
			continue
		}
		go m.handleUnixConnection(conn)
	}
}

func (m *Module) handleUnixConnection(conn *net.UnixConn) {
	var msg string
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)
	dec.Decode(&msg)
	enc.Encode(fmt.Sprintf("your message was: %s", msg))
}

func (m *Module) Shutdown() {
	m.ServerSocket.Close()
}

func (m *Module) RegisterActivity(name string) *activity.Activity {
	a := new(activity.Activity)
	a.Name = name
	m.Activities = append(m.Activities, a)
	return a
}

func Create(name string) *Module {
	m := new(Module)
	m.Name = name
	return m
}
