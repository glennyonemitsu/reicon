package main

import (
	"encoding/gob"
	"fmt"
	"github.com/glennyonemitsu/reicon/server"
	"net"
	"os"
)

func main() {
	var args []string
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("cannot get pwd")
		os.Exit(1)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "start_server" {
		server := server.CreateServer(pwd)
		server.Run()
	} else {
		addr := server.GetUnixAddress(pwd)
		conn, err := net.DialUnix("unix", nil, addr)
		if err != nil {
			fmt.Println("cannot establish connection to server")
			os.Exit(1)
			return
		}
		if len(os.Args) > 1 {
			args = os.Args[1:]
		}

		enc := gob.NewEncoder(conn)
		enc.Encode(args)

	}

}
