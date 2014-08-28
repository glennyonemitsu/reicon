package main

import (
	"fmt"
	"github.com/glennyonemitsu/reicon/server"
	"os"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("cannot get pwd")
		os.Exit(1)
	}
	server := server.CreateServer(pwd)
	server.Run()

}
